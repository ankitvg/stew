## 2026-05-01T02:08:06Z — Implement ledger cat command

**Prompt:** PLEASE IMPLEMENT THIS PLAN: Add `stew ledger cat`.

Implemented `stew ledger cat <ledger>` as a read-only command for printing raw ledger markdown.

Changes:
- Added `internal/stewledgercat` to validate ledger names, require a matching `.stew/<ledger>.spec.md`, reject missing or directory ledger paths, and read `.stew/<ledger>.md` unchanged.
- Wired `cat` into the existing `stew ledger` command group with `--path`, `stew ledger cat --help`, and `stew help ledger cat` support.
- Updated ledger help and README docs to mention reading ledgers and shell filtering.
- Added package and CLI tests for exact output, `--path`, help text, unknown ledgers, missing files, directory paths, and invalid ledger names.

Validation:
- `go test ./...`
- `go build ./...`
- `go run ./cmd/stew ledger cat iterations | sed -n '1,20p'`
- `go run ./cmd/stew ledger cat iterations | grep "Prompt" | sed -n '1,5p'`

---
