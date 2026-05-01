# Plans

<!-- Managed by stew -->

## 2026-04-27T00:05:22Z — Implement ledger creation command

**Prompt:** lets plan it then

<proposed_plan>
# Implement `stew ledger new`

## Summary
Add `stew ledger new <name>` in `/Users/ankitgupta/Documents/stewreads/stew` as the creation-side primitive for custom ledgers. The command will create `.stew/<name>.md` and `.stew/<name>.spec.md`, preserve filesystem discovery, and avoid registry/config changes.

Baseline today: `go test ./...` passes. Existing uncommitted changes documenting custom ledgers must be preserved and expanded, not reverted.

## Key Changes
- Add a new `internal/stewledger` package for ledger creation:
  - Require an existing `.stew/` directory; if missing, error with “run `stew init` first” guidance.
  - Validate names as lowercase ASCII alphanumeric plus single hyphens, 1-40 chars, no leading/trailing hyphen.
  - Reject reserved names: `stew`, `config`, `stews`, `ledger`, `ledgers`, `spec`, `specs`, plus empty or dot-prefixed names.
  - Preflight both target files and refuse to clobber if `.stew/<name>.md` or `.stew/<name>.spec.md` already exists.

- Add Cobra command surface:
  - `stew ledger new <name> [--description TEXT] [--threshold TEXT] [--path DIR] [--quiet]`
  - `--path` defaults to `.`, matching `init`, `append`, and `full-spec`.
  - `--description` and `--threshold` are optional; omitted values render conspicuous TODO guidance in the generated spec.
  - Non-quiet output should list the two created relative paths.

- Generated files:
  - `.stew/<name>.md`: H1 title derived from the ledger name plus `<!-- Managed by stew -->`, no entries.
  - `.stew/<name>.spec.md`: H1 spec title, one-line description/TODO, `## Body`, and `## When to append`.
  - `stew append <name>` and `stew full-spec --path <dir>` must work immediately after creation.

- Update existing docs/templates:
  - Expand `renderStewSpec()` and current `.stew/stew.spec.md` from the existing custom-ledger sentence into the locked ledger model: name, companion spec, shared format, append-only semantics, newest-at-bottom ordering, H2 timestamp entry boundary, `**Prompt:**` attribution, and implicit repo affiliation.
  - Update `renderManagedBlock()` and current `AGENTS.md` to mention `stew ledger new` for custom ledgers.
  - Append implementation notes to `.stew/iterations.md`; append a `.stew/decisions.md` entry because this establishes a durable first-class ledger primitive.

## Test Plan
- Add unit tests for `internal/stewledger` covering:
  - Successful creation with explicit description and threshold.
  - TODO fallback rendering when metadata flags are omitted.
  - Invalid/reserved name rejection.
  - Missing `.stew/` rejection.
  - Existing ledger/spec file rejection without clobbering.

- Add CLI/help tests covering:
  - `stew ledger new plans --path <tmp> --description ... --threshold ...`
  - `--quiet` suppressing stdout.
  - Root/ledger help documenting the new command and preserving existing command guidance.

- Update init/template tests for:
  - AGENTS managed block mentioning `stew ledger new`.
  - Stew spec documenting the canonical ledger model.

- Verify with:
  - `go test ./...`
  - `go build ./...`
  - `make pre-release`
  - Build a local temp binary, run `stew init` in a temp directory, run `stew ledger new plans ...`, then confirm `stew append plans ...` and `stew full-spec` work.

## Assumptions
- Scope is command implementation plus docs/tests only; no real `plans` ledger will be created in another repo during this pass.
- No `ledger list`, `ledger remove`, `ledger rename`, or `stew tail` is included.
- Existing append behavior remains compatible; strict new-name validation applies to `ledger new`.
</proposed_plan>

---

## 2026-05-01T01:44:14Z — Plan ledgers command

**Prompt:** Lets plan this out for stew: add a clean `stew ledgers` command listing ledger names, locations, and descriptions.

<proposed_plan>
# Add `stew ledgers`

## Summary
Add a read-only top-level command:

```sh
stew ledgers
```

It lists available writable ledgers discovered from `.stew/*.spec.md`, excluding `.stew/stew.spec.md`.

Output will be aligned, headerless rows:

```txt
decisions   .stew/decisions.md   Durable architectural or product decisions worth documenting. ...
iterations  .stew/iterations.md  A per-prompt work log. One entry per prompt that produced meaningful work. Append-only.
plans       .stew/plans.md       Complete AI-agent implementation plans, preserved for later review after they have been proposed, confirmed, or executed.
```

## Interface Changes
- Add CLI command `stew ledgers [--path <dir>]`.
- Add help support for `stew ledgers --help` and `stew help ledgers`.
- Update root help’s available command list with `ledgers    List available ledgers`.
- Do not add `stew ledger list`; keep `stew ledger new` as the creation command and `stew ledgers` as the discovery command.

## Implementation Changes
- Add a small read-only discovery package, likely `internal/stewledgers`, with:
  - `List(Options) (Result, error)`
  - `Options{TargetDir string}`
  - `Ledger{Name string, LedgerPath string, Description string}`
- Discover ledgers by reading `.stew/`, selecting non-directory files ending in `.spec.md`, excluding `stew.spec.md`, and deriving each ledger path as `.stew/<name>.md`.
- Sort ledgers alphabetically by name for deterministic output.
- Extract descriptions from the first non-empty prose paragraph after the spec H1, flattening wrapped lines into one line.
- Add `internal/cli/ledgers.go` for command parsing and tab-aligned output.
- Update README command docs to mention `stew ledgers`.

## Test Plan
- Add package tests for ledger discovery:
  - excludes `.stew/stew.spec.md`
  - sorts by ledger name
  - derives `.stew/<name>.md`
  - extracts and flattens the first description paragraph
  - errors when `.stew/` is missing
  - errors when no writable ledger specs exist
- Add CLI tests for:
  - `stew ledgers`
  - `stew ledgers --path <tmp repo>`
  - `stew ledgers --help`
  - root/help output includes the new command
- Verify with:
  - `go test ./...`
  - `go build ./...`
  - `go run ./cmd/stew ledgers`

## Assumptions
- “Available ledger” means a ledger with a matching `.stew/<name>.spec.md`; the markdown ledger file path is shown even if the file is not currently present.
- The base `.stew/stew.spec.md` is not listed because `stew` is reserved and not appendable.
- The output is intended for humans first, so aligned plain text is preferred over JSON or machine-only formatting for this iteration.
</proposed_plan>

---
