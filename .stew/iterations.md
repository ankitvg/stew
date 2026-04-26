# Iterations

<!-- Managed by stew -->

## 2026-04-26T18:46:53Z — Implement full-spec canonical flow

**Prompt:** Add to the existing init behavior by creating `.stew/stew.spec.md`, updating the AGENTS managed block to use `stew full-spec` as the canonical entry point, and making `stew full-spec` auto-include custom ledger specs.

Implemented the init flow to create six `.stew` files by adding `.stew/stew.spec.md`, updated the managed AGENTS block text to the new canonical copy, and added a new `stew full-spec` command that concatenates `.stew/*.spec.md` deterministically with `stew.spec.md` first. Added tests for new init expectations and full-spec ordering/error paths, then validated with `go test ./... && go build ./...`, `make pre-release`, plus manual smoke tests in both repo and temp non-git directories including a custom `security-audits.spec.md` file.

---

## 2026-04-26T18:52:06Z — Implement append command

**Prompt:** PLEASE IMPLEMENT THIS PLAN: Implement `stew append` as the CLI writing primitive, then commit and append logs for this session using the new local build.

Implemented `stew append` as the primary writing primitive for `.stew/` ledgers. Added a testable `internal/stewappend` package for target resolution, ledger validation, body-source handling, timestamped entry rendering, and normalized append spacing. Added the Cobra command with `--prompt`, `--summary`, `-m`, `-F`, default piped stdin support, TTY protection, and success output. Registered the command at the root and added unit/CLI tests covering formatting, validation, body source conflicts, stdin/file/message input, and required flags.

Validated with `go test ./...`, `go build ./...`, and a fresh local binary built to `/tmp/stew-local`, which was used to append this entry.

---
