## 2026-04-26T23:22:50Z — Remove ledger path flag

**Prompt:** i dont see a need to support custom file paths

Removed the user-facing `--path` flag from `stew ledger new` so ledger creation always targets the current working repository. Kept the internal `TargetDir` option for direct package tests, but the Cobra command now exposes only the ledger name plus spec metadata flags.

Updated CLI tests to chdir into temp repositories instead of passing `--path`, and added a help assertion that `ledger new` does not expose `--path`. Validation: `go test ./...`, `go build ./...`, `make pre-release`, and a temp-binary smoke test from inside the target repo all passed.

---
