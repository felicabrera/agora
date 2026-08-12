# Arquitectura de ÁGORA

Este documento describe cómo está organizado el repositorio y por qué, para que las
próximas piezas se agreguen en el lugar correcto. No describe el protocolo completo: eso
está en el documento de tesis.

## Los tres repositorios

| Repositorio | Rol | Lenguaje |
|---|---|---|
| **ÁGORA** (este) | El sistema de votación: cifra, prueba, agrega y descifra | Go |
| **[FARO](https://github.com/felicabrera/faro)** | Log público append-only sobre Tessera, sitio y CLI de auditoría | Go + Next.js |
| **[ÁBACO](https://github.com/felicabrera/abaco)** | Instrumento de medición del costo criptográfico | Go |

La dirección de las dependencias es deliberada y no debe invertirse:

```
ÁBACO  ──importa──>  ÁGORA/crypto
ÁGORA  ──cliente de──>  FARO   (HTTP, verificando todo lo que lee)
FARO   ──no conoce──>  ÁGORA
```

FARO no sabe nada de papeletas: almacena bytes opacos y prueba que no los alteró. Que esos
bytes sean una papeleta válida es una conclusión de ÁGORA. Mantener esa frontera es lo que
permite que FARO sea auditado por separado.

## `crypto/` es público, `internal/` no

`crypto/` es la superficie auditable. Está bajo AGPL-3.0 y es importable por cualquiera,
porque el argumento del proyecto es que un tercero pueda verificar un escrutinio sin
ejecutar nuestro servidor. Todo lo que se agregue ahí es, en la práctica, una promesa de
estabilidad de API.

`internal/` es pegamento del servicio: configuración, almacenamiento, HTTP. El compilador
de Go impide que se importe desde afuera, que es exactamente lo que queremos.

Regla práctica: **si un auditor lo necesita para recomputar un resultado, va en `crypto/`.
Si no, va en `internal/`.**

## Capas

```
cmd/agora-api      cmd/agorakit
      │                 │
      ├─── internal/    │        configuración, HTTP, almacenamiento, cliente FARO
      │                 │
      └────── crypto/ ──┘        grupo, ElGamal, ZKP, umbral
```

Las capas superiores dependen de las inferiores y nunca al revés. `crypto/` no importa
nada de `internal/` y no hace E/S: no abre sockets, no lee archivos y no registra logs.
Eso lo hace testeable de forma determinista y auditable sin leer el resto del sistema.

## Convenciones

Heredadas de ÁBACO, para que los tres repositorios se lean como un solo proyecto:

- **Comentarios de código en inglés**, siempre. La documentación de usuario (README,
  `docs/`) en español.
- Doc comment de paquete extenso, con un párrafo de justificación y citas académicas en
  formato autor/venue/año.
- `fmt.Errorf` con prefijo `paquete: ` en minúscula y envoltura con `%w`. Sin errores
  centinela ni tipos de error propios, salvo donde una capa superior necesite distinguir
  casos (y entonces se documenta por qué).
- Constructores `New<Tipo>`. Interfaces solo donde se ganan el lugar.
- Tests en el mismo paquete, solo con la biblioteca estándar, sin testify. Bucles
  exhaustivos sobre entradas chicas antes que tablas de casos. Vectores conocidos inline.
- `panic` reservado para invariantes documentadamente inalcanzables.

Dos divergencias deliberadas respecto de ÁBACO, porque aquí hay servicios de red y allá
una CLI: se usa `context.Context` en todo camino de E/S, y logging estructurado con
`log/slog`.

## Dependencias

Hoy hay **una sola dependencia directa**: `github.com/gtank/ristretto255`. Eso es
intencional. La superficie de dependencias de un sistema de votación forma parte de su
modelo de amenazas, y cada módulo agregado es código que hay que confiar o auditar. Antes
de agregar uno, la pregunta es si la biblioteca estándar alcanza.

Cuando llegue la integración con FARO se sumará `github.com/transparency-dev/tessera` como
cliente de lectura, que arrastra un grafo grande; conviene aislarlo detrás de
`internal/faro` para que el núcleo criptográfico siga siendo verificable por sí solo.

## Qué falta

En orden aproximado de trabajo, siguiendo el cronograma de la tesis:

1. `crypto/ballot` — la papeleta como unidad: codificar una selección, cifrarla, probarla,
   verificarla y serializarla. La serialización define la entrada del log, así que es
   consenso: dos implementaciones que no coincidan no coinciden sobre qué se votó.
2. `crypto/tally` — agregación homomórfica en streaming y lectura del exponente por
   baby-step/giant-step.
3. `internal/election` — manifiesto público de la elección y ceremonia de trustees.
4. `internal/token` — tokens de un solo uso, desvinculados de la identidad del votante.
5. `internal/faro` — cliente del log, verificando checkpoints y pruebas de inclusión.
6. `internal/storage` — padrón y estado de tokens en PostgreSQL.
7. `internal/httpapi` — emisión de voto y consulta de elección.
