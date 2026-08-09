# testr

**Persistent Eidos testing operator.** Learns how each product is *proven*, records test attempts, and keeps the test frontier concrete.

```text
ship-forge  →  shipr     (how this product ships)
test-forge  →  testr     (how this product is tested)   ← you are here
```

test-forge (methods/MCP) is **retired**. **testr** is the operator + per-repo memory.

## Install

```bash
git clone https://github.com/eidos-agi/testr.git
cd testr
pip install -e .
testr --version
```

## Per-repo file (the point)

```text
.testr/
  product-test-model.json   # how this product is tested
  test-attempts/            # ledger of runs
```

`.testr/` is gitignored by default (local operator memory), same idea as shipr’s `.shipr/`.

```bash
testr model --project /path/to/product --write --json
testr frontier --project /path/to/product --json
testr attempt --project /path/to/product \
  --goal "full suite" \
  --status passed \
  --proof "python -m pytest -q" \
  --json
```

## Relation to shipr

When a product has a testr model, **shipr `proof_commands` should align with `test_commands`**.  
Ship only after testr frontier says the gate you care about passed.

Merge plan for shipping methods: [shipr SHIP-FORGE-MERGE](https://github.com/eidos-agi/shipr/blob/main/docs/SHIP-FORGE-MERGE.md).

## Methods source

Historical testing forge: [test-forge](https://github.com/eidos-agi/test-forge) (retired).  
Best of its playbooks will absorb into testr docs/methods over time.

## License

MIT — Eidos AGI
