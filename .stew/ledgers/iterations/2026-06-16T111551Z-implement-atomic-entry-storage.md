## 2026-06-16T11:15:51Z — Implement atomic entry storage

**Prompt:** PLEASE IMPLEMENT THIS PLAN: Atomic Ledger Entry Storage

Implemented atomic ledger entry storage across the CLI. Added shared entry rendering/parsing/filename helpers, changed append to create one timestamped markdown file per entry, changed cat/tail/all-ledger reads to consume `.stew/ledgers/<ledger>/`, updated init and ledger-new to create storage directories, and added `stew migrate atomic-entries` with dry-run support.

Migrated this repo's own `decisions`, `iterations`, and `plans` ledgers from monolithic `.stew/<ledger>.md` files into `.stew/ledgers/<ledger>/`, removing the legacy files after verification. Updated README, AGENTS managed wording, the Stew spec, generated templates, CLI help, and tests.

Validation: `go test ./...`, `make build`, `./dist/stew migrate atomic-entries --dry-run`, `./dist/stew migrate atomic-entries`, and `./dist/stew ledger tail --all --limit 1 --json`.

---
