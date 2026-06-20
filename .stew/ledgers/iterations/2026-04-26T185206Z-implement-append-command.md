## 2026-04-26T18:52:06Z — Implement append command

**Prompt:** PLEASE IMPLEMENT THIS PLAN: Implement `stew append` as the CLI writing primitive, then commit and append logs for this session using the new local build.

Implemented `stew append` as the primary writing primitive for `.stew/` ledgers. Added a testable `internal/stewappend` package for target resolution, ledger validation, body-source handling, timestamped entry rendering, and normalized append spacing. Added the Cobra command with `--prompt`, `--summary`, `-m`, `-F`, default piped stdin support, TTY protection, and success output. Registered the command at the root and added unit/CLI tests covering formatting, validation, body source conflicts, stdin/file/message input, and required flags.

Validated with `go test ./...`, `go build ./...`, and a fresh local binary built to `/tmp/stew-local`, which was used to append this entry.

---
