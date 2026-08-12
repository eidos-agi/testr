# testr

> **AI testing config + proof memory** — not a test runner.  
> Sibling ship config: [shipr](https://github.com/eidos-agi/shipr).

## What this is

**`testr` is the kickstart keyword for how a product is proven.**

It stores **how-this-product-is-tested** config and a proof attempt ledger so the
next agent does not invent full-suite ceremony for every change.

It does **not** execute tests. The agent:

1. **Loads the committed model** (source of truth)
2. Runs **path-relevant** `test_commands` itself
3. Records results with `testr attempt`

```text
shipr  →  how this product ships   (.shipr/)
testr  →  how this product is tested  (.testr/)   ← you are here
```

When both exist, **shipr prefers committed shipr proofs**; if shipr proofs are
empty, it can absorb testr `test_commands`.

## 2027 kickstart contract (read this)

1. **Committed model wins.** If `.testr/product-test-model.json` exists,
   `testr model` **loads it**. It does **not** re-detect over product policy.
2. **Detection is greenfield only.** No file → starter detect (often too heavy,
   e.g. full `npm test`). Hand-edit path-relevant commands, commit, done.
3. **`--write` will not clobber.** Existing model → error unless `--force`.
4. **Lists are menus, not gates.** Prefer path-relevant commands for *this* change.
5. **Post-deploy checks** (e.g. production health curl) are not merge-time tests.

## Install (Go — canonical)

```bash
go install github.com/eidos-agi/testr/cmd/testr@latest
# or from a clone:
git clone https://github.com/eidos-agi/testr.git && cd testr
go install ./cmd/testr
testr version   # 0.3.0+
```

Legacy Python under `src/testr/` is transitional.

## Use (agent workflow)

```bash
# 1. Load THIS product's test config (committed file if present)
testr model --project /path/to/product --json

# 2. YOU run only path-relevant test_commands

# 3. Record outcome (ledger only)
testr attempt --project /path/to/product \
  --goal "relevant local proofs for this change" \
  --status passed \
  --proof "git diff --check" \
  --json

testr frontier --project /path/to/product --json
```

With shipr:

```bash
testr model --project . --json
shipr model --project . --json
# greenfield only:
testr model --project . --write --json
shipr model --project . --write --json
```

### Writing / replacing models

```bash
testr model --project . --write --json          # create if missing only
testr model --project . --write --force --json  # replace committed model (rare)
```

## Durable state

```text
.testr/
  product-test-model.json   # how this product is proven — **committed SOURCE OF TRUTH**
  test-attempts/            # ledger — **committed**
.shipr/
  product-release-model.json
```

`.testr/` and `.shipr/` are **product config, not gitignored**.

## Relation to shipr

| Tool | File | AI uses it for |
|------|------|----------------|
| testr | `.testr/product-test-model.json` | which commands prove the product |
| shipr | `.shipr/product-release-model.json` | channels, gates, rollback + proofs |

Product-specific behavior lives in those committed files — not in a parallel
custom ship system beside them.

## Agent skill

[`skills/use-testr/SKILL.md`](skills/use-testr/SKILL.md)

## License

MIT — Eidos AGI
