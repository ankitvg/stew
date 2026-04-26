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

## 2026-04-26T20:22:34Z — Dogfood ledgers and document primitive

**Prompt:** Dogfood on stewreads-docs. Migrate existing iterations.md and decisions.md into .stew/ via stew append. Also add one sentence to stew.spec.md explicitly naming custom ledgers as a primitive.

Dogfooded Stew on `stew-docs` by building a local `/tmp/stew-dogfood` binary, running `stew init`, and migrating the existing real legacy entries from `docs/iterations.md` and `docs/decisions.md` into `.stew/iterations.md` and `.stew/decisions.md` through `stew append` stdin bodies. This validated the append format against real project entries rather than synthetic examples.

Made the custom-ledger primitive explicit by adding the sentence to the Stew spec template and current `.stew/stew.spec.md`: to add a custom ledger, create `.stew/<name>.md` and `.stew/<name>.spec.md`; Stew commands discover ledgers by their spec file. Added a test to lock that generated spec guidance.

---

## 2026-04-26T23:14:23Z — Fix ledger help assertion

**Prompt:** PLEASE IMPLEMENT THIS PLAN: Implement `stew ledger new`.

Initial `go test ./...` failed in `TestLedgerNewHelpDocumentsCreationContract` because Cobra wrapped the help text between `if` and `omitted`, making a single substring assertion too brittle. The implementation behavior was intact; the next step is to assert smaller stable fragments of the same help contract.

---

## 2026-04-26T23:15:28Z — Implement ledger creation primitive

**Prompt:** PLEASE IMPLEMENT THIS PLAN: Implement `stew ledger new`.

Implemented `stew ledger new` as the custom-ledger creation primitive. Added `internal/stewledger` for strict new-ledger name validation, reserved names, `.stew/` preflight checks, no-clobber file creation, and generated ledger/spec rendering with TODO fallbacks for omitted metadata.

Added the Cobra `ledger` command with `new <name>`, `--path`, `--description`, `--threshold`, and `--quiet`; registered it at the root and updated help examples. Expanded the Stew base spec and managed AGENTS block to document first-class ledger creation and the canonical ledger property model.

Validation: `go test ./...`, `go build ./...`, `make pre-release`, and a temp-binary smoke test covering `stew init`, `stew ledger new plans`, `stew append plans`, and `stew full-spec` all passed.

---
