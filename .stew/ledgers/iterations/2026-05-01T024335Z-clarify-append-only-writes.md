## 2026-05-01T02:43:35Z — Clarify append-only writes

**Prompt:** Clarify in the Stew workflow/spec that ledger paths are for discovery and reading, not direct writes.

Clarified the Stew contract so ledger paths shown by `stew ledgers` cannot be mistaken for direct write targets.

Changes:
- Added a Core Model write-path bullet to `.stew/stew.spec.md` and `renderStewSpec()` saying ledger file paths are for discovery and reading, and entries must be created only with `stew append <ledger>`.
- Added an explicit workflow warning not to update ledger markdown files directly and to run `stew append <ledger> --help` before writing.
- Updated init/template tests to require both statements in generated specs.

Validation:
- `go test ./...`
- `go build ./...`
- `go run ./cmd/stew full-spec | sed -n '1,55p'`

---
