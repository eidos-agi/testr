# testr

> **AI testing config + proof memory** — not a test runner.  
> Succession from retired [test-forge](https://github.com/eidos-agi/test-forge).  
> Sibling ship config: [shipr](https://github.com/eidos-agi/shipr).

## What this is

**testr tells AI agents how a product is proven**, and stores that config so the
next session can run the same gates.

It is **not** an app that executes tests. The agent reads
`.testr/product-test-model.json`, runs `test_commands` itself, then records
results with `testr attempt`.

```text
shipr  →  how this product ships   (.shipr/)
testr  →  how this product is tested  (.testr/)   ← you are here
```

When both exist, **shipr prefers testr `test_commands` as `proof_commands`**.

## Install (Go — canonical)

```bash
go install github.com/eidos-agi/testr/cmd/testr@latest
# or from a clone:
git clone https://github.com/eidos-agi/testr.git && cd testr
go install ./cmd/testr
testr version
```

Legacy Python under `src/testr/` is transitional. Prefer the Go binary.

## Use (agent workflow)

```bash
# Materialize test config for AI
testr model --project /path/to/product --write --json

# AI runs the listed test_commands itself

# Record outcome (ledger only)
testr attempt --project /path/to/product \
  --goal "full suite" \
  --status passed \
  --proof "go test ./..." \
  --json

testr frontier --project /path/to/product --json
```

With shipr on the same product:

```bash
testr model --project . --write --json
shipr model --project . --write --json   # absorbs testr proofs
```

## Durable state

```text
.testr/
  product-test-model.json   # how this product is tested (AI how-to)
  test-attempts/            # ledger of proof runs
```

`.testr/` is gitignored by default (local config/memory), same idea as shipr.

## Relation to shipr

| Tool | File | AI uses it for |
|------|------|----------------|
| testr | `.testr/product-test-model.json` | which commands prove the product |
| shipr | `.shipr/product-release-model.json` | channels, gates, rollback + same proofs |

`testr frontier` surfaces `related_shipr` when a release model is present.
Ship only after the AI has run the gates you care about and recorded them.

## OPF

Canonical product graph: [`docs/opf/`](docs/opf/) — validate with `python3 -m opf.validate docs/opf`.

## Methods source

Historical testing forge: [test-forge](https://github.com/eidos-agi/test-forge) (retired).  
Playbooks absorb into testr docs over time: `docs/TEST-FORGE-SUCCESSION.md`.

## License

MIT — Eidos AGI
