## 2026-04-26T23:15:28Z — Implement ledger creation primitive

**Prompt:** PLEASE IMPLEMENT THIS PLAN: Implement `stew ledger new`.

Implemented `stew ledger new` as the custom-ledger creation primitive. Added `internal/stewledger` for strict new-ledger name validation, reserved names, `.stew/` preflight checks, no-clobber file creation, and generated ledger/spec rendering with TODO fallbacks for omitted metadata.

Added the Cobra `ledger` command with `new <name>`, `--path`, `--description`, `--threshold`, and `--quiet`; registered it at the root and updated help examples. Expanded the Stew base spec and managed AGENTS block to document first-class ledger creation and the canonical ledger property model.

Validation: `go test ./...`, `go build ./...`, `make pre-release`, and a temp-binary smoke test covering `stew init`, `stew ledger new plans`, `stew append plans`, and `stew full-spec` all passed.

---
