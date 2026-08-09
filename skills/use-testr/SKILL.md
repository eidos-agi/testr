---
name: use-testr
description: Use when the user wants product test gates, test memory, prove green, or how a repo is tested over time.
---

# Use testr

```bash
testr model --project . --write --json
testr frontier --project . --json
testr attempt --project . --goal "…" --status passed --proof "…" --json
```

Align **shipr** `proof_commands` with testr `test_commands` before `shipr attempt --status shipped`.

Methods succession: `docs/TEST-FORGE-SUCCESSION.md`. Sibling: shipr.
