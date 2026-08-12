package main

import (
	"crypto/rand"
	"encoding/hex"
	"flag"
	"fmt"

	"github.com/felicabrera/agora/crypto/elgamal"
	"github.com/felicabrera/agora/crypto/group"
	"github.com/felicabrera/agora/crypto/threshold"
)

// runKeygen deals a (t, n) threshold election key and prints it.
//
// Printing the shares to stdout is correct for a scaffold and wrong for an
// election: the point of a threshold is that the shares live with different
// people on different machines, and a command that emits all of them has already
// defeated it. Distribution is part of the trustee ceremony, which is MVP work.
func runKeygen(args []string) error {
	fs := flag.NewFlagSet("keygen", flag.ExitOnError)
	authorities := fs.Int("authorities", 5, "number of key-holding authorities (n)")
	thr := fs.Int("threshold", 3, "how many authorities are needed to decrypt (t)")
	showShares := fs.Bool("unsafe-print-shares", false,
		"print every authority's secret share to stdout (development only; destroys the threshold)")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if fs.NArg() > 0 {
		return fmt.Errorf("keygen: unexpected argument %q", fs.Arg(0))
	}
	if *thr < 1 || *thr > *authorities {
		return fmt.Errorf("keygen: need 1 <= threshold <= authorities, got %d of %d", *thr, *authorities)
	}

	g := group.NewRistretto255()
	x, shares, err := threshold.Deal(g, *authorities, *thr, rand.Reader)
	if err != nil {
		return fmt.Errorf("keygen: dealing the key: %w", err)
	}
	pk := elgamal.PublicKeyFrom(g, x)

	fmt.Printf("group:      %s\n", g.Name())
	fmt.Printf("threshold:  %d of %d\n", *thr, *authorities)
	fmt.Printf("public key: %s\n", hex.EncodeToString(pk.Y.Bytes()))
	fmt.Printf("shares:     %d generated\n", len(shares))

	if !*showShares {
		fmt.Println()
		fmt.Println("The secret shares were NOT printed. Any t of them reconstruct the decryption")
		fmt.Println("key, so collecting them in one terminal, one scrollback buffer or one shell")
		fmt.Println("history is the same as not having a threshold at all. Distributing them to")
		fmt.Println("their holders is the trustee ceremony, which this command does not implement.")
		fmt.Println()
		fmt.Println("For local development only, re-run with --unsafe-print-shares.")
		return nil
	}

	fmt.Println()
	fmt.Println("!! DEVELOPMENT ONLY. Every share below is secret. Together they ARE the")
	fmt.Println("!! decryption key. Never run this against an election anyone relies on.")
	for _, s := range shares {
		fmt.Printf("  %2d: %s\n", s.Index, hex.EncodeToString(s.Value.Bytes()))
	}
	return nil
}
