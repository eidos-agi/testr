# test-forge → testr succession

| | test-forge | testr |
|--|------------|-------|
| Status | Retired 2026-08-03 | Active |
| Form | MCP / playbooks / probe | Go CLI + `.testr/` config + attempts |
| Job | “Test this” methods | **AI how-to-prove config** + proof ledger |

## Product identity

testr does **not** run tests. It stores repeatable test config so AI agents run
the same gates every time. Sibling **shipr** absorbs `test_commands` as
`proof_commands` when both models exist.

## Absorb later

- Playbooks and structured evidence patterns from test-forge / probe_forge
- Optional “suggest command” helpers — still not a runner

## Do not

- Re-enable test-forge MCP as the brand
- Build a test *runner* inside testr — compose: AI runs commands; testr/shipr store config

## forge-forge registry

- `test-forge` entry: `status: retired`, `successor: testr`
- `testr` entry: active `type: tool`, invocation `cli`, Go install
- `shipr` is the shipping sibling; both configs are committed (not gitignored)
