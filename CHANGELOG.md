# Changelog

## 0.3.0

- **Committed model wins:** `testr model` loads `.testr/product-test-model.json` when present (`model_source: committed`).
- Detection is greenfield-only; never invent full-suite ceremony over a committed file.
- `--write` refuses to overwrite without `--force`.
- CLI help, README, and `use-testr` skill document the 2027 kickstart contract.
