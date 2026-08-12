---
name: use-testr
description: >
  Use when the user wants product test gates, test memory, prove green, or how
  a repo is tested. testr stores AI config + proof memory — it does not run
  tests itself. Prefer committed product models over auto-detect.
---

# Use testr

testr is **AI testing config + proof memory**. It does **not** execute tests.
You load the product model, run **path-relevant** `test_commands` yourself, then
record results.

## 2027 kickstart (do this order)

```bash
# 1. LOAD committed model — never lead with blind --write
testr model --project . --json
# Look for model_source: "committed". That file is product policy.

# 2. YOU run only path-relevant test_commands for THIS change

# 3. Record
testr attempt --project . \
  --goal "…" \
  --status planned|running|passed|failed|blocked|skipped \
  --proof "…" \
  --json
```

### Hard rules

1. **Committed `.testr/product-test-model.json` wins.** If present, obey it.
2. **Do not** run `testr model --write` when a model already exists unless the
   human asked to regenerate (`--force`).
3. **Lists are menus.** Do not run the entire `test_commands` array on every
   change; pick what the diff touches. Notes like `test_command_notes` matter.
4. **Detection is greenfield only.** Full `npm test` from detect is a starter
   guess, not a product mandate.
5. Health / production curls are **post-deploy**, not merge-time gates (see shipr).

## Align with shipr

```bash
testr model --project . --json
shipr model --project . --json
```

Before `shipr attempt --status shipped`, record the proofs you actually ran.

## Install

```bash
go install github.com/eidos-agi/testr/cmd/testr@latest
testr version   # expect 0.3.0+
```

## Tracked config

`.testr/` and `.shipr/` are committed product config. Tools refuse to overwrite
existing models without `--force`.
