# Política de Seguridad — ÁGORA

ÁGORA es un sistema de voto electrónico auditable y criptográficamente
verificable. El secreto del voto y la integridad del escrutinio son la razón
de ser del sistema: cualquier falla que permita vincular un voto con su
votante, alterar el conteo o degradar las garantías criptográficas se trata
como crítica.

## Reportar una vulnerabilidad

Si encontrás una vulnerabilidad de seguridad, una debilidad criptográfica o un
problema de integridad, reportalo **de forma privada** usando la función de
reporte privado de GitHub:

> Pestaña **Security** → **Report a vulnerability**

No abras un issue público ni lo divulgues en redes antes de que exista un
arreglo disponible.

- **Acuse de recibo:** dentro de las 72 horas.
- **Plan de resolución:** dentro de los 7 días desde el acuse de recibo.
- **Divulgación coordinada:** se espera que nos permitas publicar el arreglo
  antes de cualquier divulgación pública. Si el reporte es válido, damos crédito
  público al reportante salvo que prefiera permanecer anónimo.

## Alcance

**Dentro de alcance**

- Implementación criptográfica: cifrado homomórfico ElGamal, pruebas de
  conocimiento cero (Schnorr, Bulletproofs), descifrado con umbral
  (threshold decryption) y generación distribuida de claves.
- Fugas de información sobre el texto plano del voto: canales laterales,
  operaciones no constant-time, fuentes de aleatoriedad débiles, validación
  faltante sobre puntos de curva o entradas criptográficas.
- Integridad del escrutinio: agregación homomórfica, verificación de pruebas,
  correspondencia entre votos emitidos y votos contados.
- Integridad del log de transparencia FARO y su verificación desde ÁGORA.
- Autenticación de votantes, autorización, manejo de sesiones y tokens.
- API del sistema de voto y configuración de infraestructura y CI/CD.

**Fuera de alcance**

- Vulnerabilidades en dependencias de terceros sin un vector de explotación
  demostrable en este proyecto (reportalas upstream).
- Denegación de servicio por volumen de tráfico.
- Hallazgos automatizados de escáneres sin análisis de impacto.
- Ingeniería social contra el equipo, la institución o votantes reales.

## Versiones soportadas

| Versión        | Soportada |
| -------------- | --------- |
| `main` (última)| ✅        |

Este es un proyecto académico en desarrollo activo (Proyecto Final de Grado,
Ingeniería en Informática, Uruguay). No hay versiones estables anteriores con
soporte.

---

# Security Policy (English)

ÁGORA is an auditable, cryptographically verifiable electronic voting system.
Ballot secrecy and tally integrity are the system's entire purpose: any flaw
that links a vote to its voter, alters the count, or weakens the cryptographic
guarantees is treated as critical.

## Reporting a Vulnerability

If you discover a security vulnerability, cryptographic weakness, or integrity
issue, please report it **privately** through GitHub's private vulnerability
reporting feature (**Security** tab → **Report a vulnerability**) rather than
opening a public issue.

- **Acknowledgement:** within 72 hours.
- **Resolution timeline:** within 7 days of acknowledgement.
- **Coordinated disclosure** is expected; please allow us to ship a fix before
  public disclosure. Valid reports are credited publicly unless you prefer to
  remain anonymous.

## Scope

**In scope:** cryptographic implementation (ElGamal homomorphic encryption,
Schnorr/Bulletproofs zero-knowledge proofs, threshold decryption, distributed
key generation); information leakage about vote plaintext (side channels,
non-constant-time operations, weak randomness, missing validation of curve
points); tally integrity; FARO transparency log integrity and its verification
from ÁGORA; voter authentication, authorization, session and token handling;
the voting system API, and infrastructure/CI-CD configuration.

**Out of scope:** third-party dependency issues with no demonstrated exploit
path in this project, volumetric denial of service, unanalyzed automated
scanner output, and social engineering.

## Supported Versions

| Version         | Supported |
| --------------- | --------- |
| `main` (latest) | ✅        |
