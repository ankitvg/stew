## 2026-05-01T01:44:35Z — Implement ledgers command

**Prompt:** PLEASE IMPLEMENT THIS PLAN: Add `stew ledgers` command listing ledger names, locations, and descriptions.

Implemented a read-only `stew ledgers` top-level command.

Changes:
- Added `internal/stewledgers` to discover writable ledgers from `.stew/*.spec.md`, exclude `.stew/stew.spec.md`, derive `.stew/<name>.md`, sort by ledger name, and flatten the first prose paragraph after each spec H1 into the description.
- Added `internal/cli/ledgers.go`, root dispatch, help support, and tab-aligned headerless output for `stew ledgers [--path <dir>]`.
- Updated README command docs to include `stew ledgers`.
- Added unit tests for ledger discovery and CLI tests for listing, `--path`, help, and root help coverage.

Validation:
- `go test ./...`
- `go build ./...`
- `go run ./cmd/stew ledgers`

---
