package threshold

import (
	"crypto/rand"
	"testing"

	"github.com/felicabrera/agora/crypto/elgamal"
	"github.com/felicabrera/agora/crypto/group"
)

// decryptWith recombines the given shares against ciphertext ct and returns the
// recovered plaintext point B - x*A.
func decryptWith(g group.Group, shares []Share, ct *elgamal.Ciphertext) group.Element {
	partials := make([]PartialDecryption, len(shares))
	for k, s := range shares {
		partials[k] = s.PartialDecrypt(ct.A)
	}
	xA, err := Combine(g, partials)
	if err != nil {
		panic("threshold: test helper combined an invalid quorum: " + err.Error())
	}
	return ct.B.Sub(xA)
}

func TestThresholdDecryptsWithQuorum(t *testing.T) {
	g := group.NewRistretto255()
	const n, thresh = 5, 3

	x, shares, err := Deal(g, n, thresh, rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	pk := elgamal.PublicKeyFrom(g, x)

	// Encrypt m=1 and recover it through a quorum of authorities.
	ct, _, err := elgamal.EncryptRandom(pk, 1, rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	got := decryptWith(g, shares[:thresh], ct)
	want := elgamal.PointFromCount(g, 1)
	if !got.Equal(want) {
		t.Fatal("quorum of t shares failed to decrypt")
	}
}

func TestThresholdAnyQuorumAgrees(t *testing.T) {
	g := group.NewRistretto255()
	const n, thresh = 5, 3
	x, shares, _ := Deal(g, n, thresh, rand.Reader)
	pk := elgamal.PublicKeyFrom(g, x)
	ct, _, _ := elgamal.EncryptRandom(pk, 1, rand.Reader)

	// Two different quorums must recover the same plaintext point.
	a := decryptWith(g, []Share{shares[0], shares[1], shares[2]}, ct)
	b := decryptWith(g, []Share{shares[1], shares[3], shares[4]}, ct)
	if !a.Equal(b) {
		t.Fatal("different quorums disagreed on the decryption")
	}
	if !a.Equal(elgamal.PointFromCount(g, 1)) {
		t.Fatal("quorum decrypted to the wrong value")
	}
}

// TestSubThresholdCannotDecrypt is the safety property: t-1 shares must not
// recover the plaintext.
func TestSubThresholdCannotDecrypt(t *testing.T) {
	g := group.NewRistretto255()
	const n, thresh = 5, 3
	x, shares, _ := Deal(g, n, thresh, rand.Reader)
	pk := elgamal.PublicKeyFrom(g, x)
	ct, _, _ := elgamal.EncryptRandom(pk, 1, rand.Reader)

	got := decryptWith(g, shares[:thresh-1], ct) // only t-1 shares
	if got.Equal(elgamal.PointFromCount(g, 1)) {
		t.Fatal("t-1 shares recovered the plaintext: threshold is broken")
	}
}

func TestCombineRecoversKeyTimesA(t *testing.T) {
	// Combine of a quorum must equal x*A exactly.
	g := group.NewRistretto255()
	const n, thresh = 7, 4
	x, shares, _ := Deal(g, n, thresh, rand.Reader)

	a, _ := g.RandomScalar(rand.Reader)
	A := g.ScalarBaseMul(a)

	partials := make([]PartialDecryption, thresh)
	for k := 0; k < thresh; k++ {
		partials[k] = shares[k].PartialDecrypt(A)
	}
	combined, err := Combine(g, partials)
	if err != nil {
		t.Fatalf("Combine: %v", err)
	}
	if !combined.Equal(A.ScalarMul(x)) {
		t.Fatal("Lagrange recombination did not yield x*A")
	}
}

// TestCombineRejectsRepeatedAuthority guards the failure mode that motivated
// validating Combine's input. Before it validated, a quorum containing the same
// authority twice produced a well-formed group element that was simply the wrong
// answer: Lagrange interpolation is undefined over repeated indices, and because
// ristretto255 defines the inverse of zero as zero, the degenerate denominator
// went unnoticed. A wrong x*A means a wrong tally, published with proofs and
// indistinguishable from a correct one.
func TestCombineRejectsRepeatedAuthority(t *testing.T) {
	g := group.NewRistretto255()
	_, shares, err := Deal(g, 5, 3, rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	A := g.ScalarBaseMul(g.ScalarFromUint64(7))

	partials := []PartialDecryption{
		shares[0].PartialDecrypt(A),
		shares[0].PartialDecrypt(A), // the same authority, twice
		shares[1].PartialDecrypt(A),
	}
	if _, err := Combine(g, partials); err == nil {
		t.Fatal("Combine accepted a quorum containing the same authority twice")
	}
}

func TestCombineRejectsMalformedQuorums(t *testing.T) {
	g := group.NewRistretto255()
	_, shares, err := Deal(g, 3, 2, rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	A := g.ScalarBaseMul(g.ScalarFromUint64(3))
	valid := shares[0].PartialDecrypt(A)

	cases := []struct {
		name     string
		partials []PartialDecryption
	}{
		{"empty", nil},
		// Index 0 evaluates the polynomial at its constant term, which is the
		// secret itself. It is never dealt, so it must never be accepted.
		{"index zero", []PartialDecryption{{Index: 0, D: A}, valid}},
		{"missing value", []PartialDecryption{{Index: 2, D: nil}, valid}},
	}
	for _, tc := range cases {
		if _, err := Combine(g, tc.partials); err == nil {
			t.Fatalf("%s: Combine = nil error, want an error", tc.name)
		}
	}
}
