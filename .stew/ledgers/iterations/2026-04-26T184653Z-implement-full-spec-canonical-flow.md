## 2026-04-26T18:46:53Z — Implement full-spec canonical flow

**Prompt:** Add to the existing init behavior by creating `.stew/stew.spec.md`, updating the AGENTS managed block to use `stew full-spec` as the canonical entry point, and making `stew full-spec` auto-include custom ledger specs.

Implemented the init flow to create six `.stew` files by adding `.stew/stew.spec.md`, updated the managed AGENTS block text to the new canonical copy, and added a new `stew full-spec` command that concatenates `.stew/*.spec.md` deterministically with `stew.spec.md` first. Added tests for new init expectations and full-spec ordering/error paths, then validated with `go test ./... && go build ./...`, `make pre-release`, plus manual smoke tests in both repo and temp non-git directories including a custom `security-audits.spec.md` file.

---
