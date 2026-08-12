# ÁGORA

**Auditable Government Open Registry Architecture** — the cryptographic core of an
auditable electronic voting system: homomorphic ElGamal, zero-knowledge ballot-validity
proofs and threshold decryption, in Go.

> **English (short version).** ÁGORA is the voting system itself: it encrypts a ballot,
> proves the ballot is valid without revealing it, aggregates ballots homomorphically and
> decrypts only the final total, under a threshold of authorities. Its companion
> [FARO](https://github.com/felicabrera/faro) is the public append-only log that makes the
> result independently auditable, and [ÁBACO](https://github.com/felicabrera/abaco)
> measures what the cryptography costs. This repository currently contains the
> cryptographic core and the service skeletons; the protocol layer is under development.
> Documentation below is in Spanish, the language of its primary audience.

---

## Estado actual

Este repositorio está en la fase inicial de desarrollo. Lo que hay hoy:

| Componente | Estado |
|---|---|
| `crypto/group` | Grupo de orden primo (ristretto255, RFC 9496) detrás de una interfaz |
| `crypto/elgamal` | ElGamal exponencial (aditivamente homomórfico) |
| `crypto/zkp` | OR-proof `{0,1}` (CDS) + prueba 1-de-C (Chaum-Pedersen), Fiat-Shamir |
| `crypto/threshold` | Shamir + descifrado por umbral con interpolación de Lagrange |
| `cmd/agora-api` | Esqueleto del servidor: configuración, logging, `/healthz`, apagado ordenado |
| `cmd/agorakit` | CLI de operador: `version`, `keygen` |

Pendiente, en orden de trabajo: papeleta y agregación (`crypto/ballot`, `crypto/tally`),
padrón y tokens de un solo uso, integración con FARO, API de elección y panel de
administración.

## Relación con ÁBACO

El núcleo criptográfico de `crypto/` proviene de
[ÁBACO](https://github.com/felicabrera/abaco), donde se escribió y se midió primero. ÁBACO
se describe a sí mismo como *"benchmark suite for the ÁGORA verifiable-voting core"*, y
hasta ahora medía una copia privada de ese núcleo porque el núcleo no existía como
artefacto separado. Vive aquí, público e importable bajo AGPL-3.0, para que:

1. ÁBACO mida el código que realmente se ejecuta, no una copia;
2. un tercero pueda verificar un escrutinio sin ejecutar el servidor de ÁGORA;
3. exista una sola implementación de ElGamal en el proyecto y no dos.

El paso de ÁBACO a importar `github.com/felicabrera/agora/crypto` es trabajo posterior, en
aquel repositorio.

## Quickstart

Requiere **Go 1.26+**.

```console
$ git clone https://github.com/felicabrera/agora && cd agora
$ make build

# Genera una clave de elección con umbral 3 de 5:
$ ./bin/agorakit keygen --authorities 5 --threshold 3

# Levanta el esqueleto de la API:
$ ./bin/agora-api
$ curl localhost:8080/healthz
```

`make help` lista el resto de los targets.

## Especificación criptográfica

Criptografía estándar y documentada; no se inventan variantes. Las referencias están
citadas en los comentarios de cada paquete.

**Grupo.** [ristretto255](https://ristretto.group/) (RFC 9496) sobre Curve25519: grupo de
orden primo, sin cofactor ni subgrupos pequeños, que es lo que ElGamal y Schnorr asumen.
Está detrás de una interfaz (`crypto/group`) para permitir un segundo backend.

**ElGamal exponencial.** Clave privada `x`, pública `Y = xG`. Cifrado de `m` con
aleatoriedad fresca `r`: `A = rG`, `B = rY + mG`. Homomórfico: `(A₁+A₂, B₁+B₂)` cifra
`m₁+m₂`. La frescura de `r` da seguridad semántica (IND-CPA), que es lo que permite que
FARO sea público. *(ElGamal 1985; Cramer-Gennaro-Schoenmakers, EUROCRYPT '97.)*

**ZKP de validez `{0,1}`.** Prueba disyuntiva de que un ciphertext cifra 0 **o** 1 sin
revelar cuál: composición OR estilo CDS de dos sigma-protocolos Chaum-Pedersen, no
interactiva por Fiat-Shamir. El hash incluye un separador de dominio, la clave pública y
el ciphertext completo; omitirlos rompe la solidez. *(Cramer-Damgård-Schoenmakers,
CRYPTO '94; Chaum-Pedersen, CRYPTO '92; Fiat-Shamir, CRYPTO '86.)*

**Prueba 1-de-C.** La OR-proof por casilla garantiza `{0,1}` en cada una; para que la
papeleta sea un voto real hace falta probar que cifra **exactamente un 1**. Se agrega una
prueba Chaum-Pedersen sobre el ciphertext agregado. Una papeleta que vota a dos candidatos
es rechazada.

**Descifrado por umbral.** Polinomio de grado `t−1` sobre `Z_q` con `x = f(0)`; shares
`(i, f(i))`. Cada autoridad publica `Dᵢ = xᵢ·A`; la interpolación de Lagrange en 0 recupera
`x·A`. Con `t−1` shares no se puede. *(Shamir, CACM 1979.)*

## Límites conocidos

Declararlos protege la honestidad del proyecto:

- **La implementación criptográfica es de referencia, no de producción.** No fue auditada
  por terceros ni endurecida contra ataques de canal lateral (timing, caché).
- **El reparto de shares usa un dealer confiable.** El generador conoce `x` durante la
  ceremonia. La variante sin dealer (Pedersen, EUROCRYPT '91) es la mejora prevista.
- **No hay todavía capa de protocolo.** Padrón, tokens, API de emisión y escrutinio están
  pendientes; lo que existe es el núcleo y los esqueletos de servicio.

## Estructura del repositorio

```
crypto/group/        # interfaz del grupo + backend ristretto255
crypto/elgamal/      # ElGamal exponencial
crypto/zkp/          # OR-proof {0,1} + prueba 1-de-C
crypto/threshold/    # Shamir + descifrado por umbral
cmd/agora-api/       # servidor de la API de elección
cmd/agorakit/        # CLI de operador
internal/config/     # configuración por entorno
internal/version/    # información de build estampada por el toolchain
```

## Tests

```console
$ make test-race     # corrección criptográfica, con race detector
$ make lint          # golangci-lint v2
```

## Seguridad

Ver [`SECURITY.md`](SECURITY.md). Para reportar una vulnerabilidad, usar el canal de
divulgación coordinada descrito allí y no un issue público.

## Licencia

GNU AGPL-3.0 — ver [`LICENSE`](LICENSE). La licencia es deliberada: un sistema de votación
cuyo argumento central es la auditabilidad no puede desplegarse como servicio sin publicar
el código que se está ejecutando.

## Referencias

ElGamal (1985) · Cramer, Damgård, Schoenmakers, *Proofs of Partial Knowledge*, CRYPTO '94 ·
Cramer, Gennaro, Schoenmakers, *A Secure and Optimally Efficient Multi-Authority Election
Scheme*, EUROCRYPT '97 · Chaum, Pedersen, CRYPTO '92 · Fiat, Shamir, CRYPTO '86 · Shamir,
*How to Share a Secret*, CACM 1979 · Pedersen, EUROCRYPT '91 ·
[RFC 9496 (ristretto255)](https://www.rfc-editor.org/rfc/rfc9496) ·
[RFC 6962 (Certificate Transparency)](https://www.rfc-editor.org/rfc/rfc6962) ·
[C2SP tlog-tiles](https://c2sp.org/tlog-tiles) ·
[Tessera](https://github.com/transparency-dev/tessera)
