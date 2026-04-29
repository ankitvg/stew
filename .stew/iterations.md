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

## 2026-04-26T23:22:50Z — Remove ledger path flag

**Prompt:** i dont see a need to support custom file paths

Removed the user-facing `--path` flag from `stew ledger new` so ledger creation always targets the current working repository. Kept the internal `TargetDir` option for direct package tests, but the Cobra command now exposes only the ledger name plus spec metadata flags.

Updated CLI tests to chdir into temp repositories instead of passing `--path`, and added a help assertion that `ledger new` does not expose `--path`. Validation: `go test ./...`, `go build ./...`, `make pre-release`, and a temp-binary smoke test from inside the target repo all passed.

---

## 2026-04-27T00:06:03Z — Dogfood plans ledger

**Prompt:** awesome, shall we test it. I want to create a plans ledger that stores all the plans that I create with AI agents. I think it could be pretty cool to see all implemented plans after a while. I also think we should try to save it in entirety, whatever the model gave us plan to confirm, just save it. So we reduce work on model's end.

Created a real `.stew/plans.md` and `.stew/plans.spec.md` in the stew repo using `stew ledger new plans` without `--path`. Tuned the generated plans spec so entries preserve the full AI-generated plan exactly as proposed, including headings, bullets, code blocks, and `<proposed_plan>` tags, with only an optional outcome note after the copied plan.

Appended the first real plans entry: the full previously confirmed `stew ledger new` implementation plan. Verified `stew full-spec` discovers the new plans spec and `go test ./...` still passes.

---

## 2026-04-28T02:34:37Z — Prepare release docs

**Prompt:** Lets skip license for a bit more? Do they rest

Expanded README with install, quick-start, command, and release-check guidance for v0.1.0. Updated Makefile so make build creates a versioned dist/stew binary with linker-injected version, commit, and date metadata, while make pre-release still runs tests and package builds first. Validated with make pre-release VERSION=v0.1.0 and ./dist/stew version.

---

## 2026-04-28T02:36:50Z — Publish first release

**Prompt:** commit and push\nand create a release on gh

Committed and published the v0.1.0 release prep. Re-ran stew help and stew full-spec, verified make pre-release VERSION=v0.1.0 passed locally, checked GitHub authentication, and confirmed no existing v0.1.0 tag or release was present before publishing.

---

## 2026-04-28T03:17:58Z — Remove Cobra dependency

**Prompt:** PLEASE IMPLEMENT THIS PLAN: Remove Cobra And Use Standard Library CLI

Replaced the Cobra-based command layer with a standard-library dispatcher using flag.FlagSet, explicit help text, and ExecuteWithIO for tests. Preserved documented commands, root version/help flags, command-specific help, and append flags after the ledger positional. Removed Cobra, pflag, and mousetrap from go.mod and deleted go.sum because no third-party dependencies remain. Updated CLI tests away from Cobra command helpers and added a regression check that root help no longer lists completion. Validated with go test ./..., make pre-release VERSION=v0.1.0, manual help/version checks, a temp-repo append using post-positional flags, and an rg check for Cobra dependency residue.

---

## 2026-04-28T03:19:59Z — Commit stdlib CLI

**Prompt:** commit and push

Committed and pushed the standard-library CLI rewrite. Before committing, re-ran stew help and stew full-spec, verified make pre-release VERSION=v0.1.0 passed, and confirmed no Cobra, pflag, mousetrap, or spf13 references remain in go.mod, internal, or cmd.

---

## 2026-04-28T03:23:02Z — Document agent-first usage

**Prompt:** The current README is actually pretty great. But before install, I want to add an appropriately titled section which basically conveys that stew will most of the time used only by ai agent. from init to append to new ledger creation. Its the AI client that can take care of all this. User doesnt have to really learn anything

Added an Agent-First Usage section before Install in README.md. The section explains that Stew is primarily meant to be operated by an AI coding agent after installation: initializing a repo, reading the full contract, appending iteration or decision entries, and creating new ledgers as needed. Validated with go test ./....

---

## 2026-04-28T03:29:31Z — Clarify agent setup loop

**Prompt:** nice but may be we should mention the entire setup. After all, stew init indeed plugs stew into AGENTS.md, so in a way all user has to do is init stew which cli agent can even do

Refined README.md's Agent-First Usage section to describe the full bootstrap flow: install the CLI, run or ask an agent to run stew init, and let the managed AGENTS.md block teach future agent sessions to run stew help, load full-spec, append entries, and create ledgers. Validated with go test ./....

---

## 2026-04-28T05:21:09Z — Commit README update

**Prompt:** commit and push

Committed and pushed the README agent-first setup update. Before committing, re-ran stew help and stew full-spec, confirmed the pending diff was limited to README.md and .stew/iterations.md, and validated with go test ./....

---

## 2026-04-29T10:09:35Z — Add Apache license

**Prompt:** You know what, lets put a APACHE 2.0 license for stew cli

Added the Apache License 2.0 text for the stew CLI and documented the license in README.md. Validated with go test ./... and go build ./....

---
