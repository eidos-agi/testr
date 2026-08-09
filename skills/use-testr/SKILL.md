---
name: use-testr
description: Use when the user wants product test gates, test memory, prove green, or how a repo is tested over time. testr stores AI config + proof memory — it does not run tests itself.
---

# Use testr

testr is **AI testing config + proof memory**. It does **not** execute tests.
You read the model, run `test_commands` yourself, then record results.

```bash
testr model --project . --write --json
# YOU run the test_commands from the model
testr attempt --project . --goal "…" --status passed --proof "…" --json
testr frontier --project . --json
```

## Align with shipr

```bash
testr model --project . --write --json
shipr model --project . --write --json   # shipr absorbs test_commands as proofs
```

Before `shipr attempt --status shipped`, the gates you care about should already
be recorded via testr (or run as shipr `proof_commands`).

## Install

```bash
go install github.com/eidos-agi/testr/cmd/testr@latest
```

Methods succession: `docs/TEST-FORGE-SUCCESSION.md`. Sibling: shipr.

## Tracked config

`.shipr/` and `.testr/` are committed product config. Tools strip ignore rules and create missing sibling models on write.
