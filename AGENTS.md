# testr — agent notes

## OPF

Canonical product graph: `docs/opf/product.json`. Validate: `python3 -m opf.validate docs/opf`.

## Product identity

testr is **AI testing config + proof memory**. It does **not** execute tests.
You (the agent) read the model, run `test_commands` yourself, then record outcomes.

## Canonical CLI

Go binary: `go install ./cmd/testr` (or `@latest`). Python `src/testr/` is legacy.

## Per-product files

```text
.testr/product-test-model.json   # how this product is proven
.testr/test-attempts/            # ledger (gitignored)
```

## Workflow

1. `testr model --project . --write --json`
2. Run each `test_commands` entry yourself
3. `testr attempt --goal "…" --status planned|running|passed|failed|blocked|skipped --proof "…" --json`
4. `testr frontier --project . --json`

## Compose with shipr

Write testr first so shipr can absorb `test_commands` as `proof_commands`:

```bash
testr model --project . --write --json
shipr model --project . --write --json
```

`related_shipr` on the model/frontier points at `.shipr/product-release-model.json` when present.
