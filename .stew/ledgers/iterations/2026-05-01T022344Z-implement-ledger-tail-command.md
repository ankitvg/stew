## 2026-05-01T02:23:44Z — Implement ledger tail command

**Prompt:** PLEASE IMPLEMENT THIS PLAN: Add `stew ledger tail`.

Implemented `stew ledger tail <ledger>` as an entry-aware read command.

Changes:
- Added `internal/stewledgertail` to reuse `stewledgercat.Run`, validate positive limits, detect canonical Stew H2 timestamp headings, and return the final N entries without the ledger preamble.
- Wired `tail` into the existing `stew ledger` command group with `--limit` defaulting to 10 and `--path` matching other read commands.
- Added `stew ledger tail --help`, `stew help ledger tail`, and ledger command help text.
- Updated README command docs with a tail example.
- Added package and CLI tests covering limit behavior, ordering, no-entry ledgers, preamble omission, path handling, help, and non-entry H2 headings inside entry bodies.

Validation:
- `go test ./...`
- `go build ./...`
- `go run ./cmd/stew ledger tail iterations --limit 5`
- `go run ./cmd/stew ledger tail --help`
- `go run ./cmd/stew help ledger tail`

---
