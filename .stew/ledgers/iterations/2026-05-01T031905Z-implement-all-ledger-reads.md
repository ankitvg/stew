## 2026-05-01T03:19:05Z — Implement all-ledger reads

**Prompt:** PLEASE IMPLEMENT THIS PLAN: Add `--all` to `stew ledger cat` and `stew ledger tail`

Implemented all-ledger read support for `stew ledger cat` and `stew ledger tail`.

Changes:
- Added `internal/stewledgerall` to discover writable ledgers, read each ledger in sorted name order, and return `{Name, Content}` sections for `cat --all` and `tail --all`.
- Wired `--all` into `stew ledger cat` and `stew ledger tail`, with clear mutual exclusion against a ledger positional.
- Preserved existing single-ledger behavior: `cat <ledger>` remains raw ledger markdown and `tail <ledger>` remains selected entries only.
- Rendered all-ledger output under command-facing `# <ledger>` headings without adding ledger file paths.
- Updated help text, README, `stew.spec.md`, and init templates so startup memory loading uses `stew ledger tail --all --limit 5`.
- Added package and CLI tests covering sorted all-ledger output, per-ledger tail limits, empty ledger sections, missing content errors, path support, mutual exclusion, missing positional behavior, and help text.

Validation:
- `go test ./...`
- `go build ./...`
- `go run ./cmd/stew ledger cat --all`
- `go run ./cmd/stew ledger tail --all --limit 5`

---
