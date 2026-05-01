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

## 2026-05-01T02:07:57Z — Plan ledger cat command

**Prompt:** I think we should implement "stew ledger cat <ledger>" next. It should print the entire ledger and support shell search via pipes.

<proposed_plan>
# Add `stew ledger cat`

## Summary
Add a read-only ledger printing command:

```sh
stew ledger cat <ledger>
```

It prints the complete raw `.stew/<ledger>.md` file to stdout so users and agents can pipe it into shell tools:

```sh
stew ledger cat iterations | grep "parser"
```

## Interface Changes
- Add `stew ledger cat <ledger> [--path <dir>]`.
- Add help support for `stew ledger cat --help` and `stew help ledger cat`.
- Update `stew ledger --help` to list:
  - `cat         Print a stew ledger`
  - `new         Create a custom stew ledger`
- Output is raw file content only, with no banner, metadata, path prefix, or extra summary line.

## Implementation Changes
- Add a small read-only package, likely `internal/stewledgercat`, with `Run(Options) (Result, error)`.
- Validate the ledger name the same way append does: require a base name, reject `stew`, whitespace, path traversal, and path separators.
- Require `.stew/<ledger>.spec.md` to exist so `cat` only works for known writable ledgers.
- Read `.stew/<ledger>.md` and return its content unchanged.
- Error clearly if the spec is missing, the ledger file is missing, or either path is a directory.
- Add CLI parsing in `internal/cli/ledger.go`, keeping `cat` under the existing `ledger` command group.

## Test Plan
- Add package tests covering:
  - prints exact ledger file content
  - rejects unknown ledger when `.stew/<ledger>.spec.md` is absent
  - rejects missing ledger file when spec exists but `.md` is absent
  - rejects reserved/path-like ledger names
  - supports `TargetDir`
- Add CLI tests covering:
  - `stew ledger cat iterations`
  - `stew ledger cat iterations --path <tmp repo>`
  - `stew ledger cat --help`
  - `stew help ledger cat`
  - `stew ledger --help` includes `cat`
- Verify with:
  - `go test ./...`
  - `go build ./...`
  - `go run ./cmd/stew ledger cat iterations`
  - `go run ./cmd/stew ledger cat iterations | grep "Prompt"`

## Assumptions
- `cat` means full raw ledger output, including the H1 and `<!-- Managed by stew -->` marker.
- Search, tail, entry limiting, and structured filters are intentionally left to shell pipes for now.
- `--path` is included because this is a read-only repository inspection command, matching `full-spec`, `append`, and `ledgers`.
</proposed_plan>

---

## 2026-05-01T02:23:35Z — Plan ledger tail command

**Prompt:** Plan entry-aware `stew ledger tail <ledger> --limit N` for recent ledger entries.

<proposed_plan>
# Add `stew ledger tail`

## Summary
Add an entry-aware read command:

```sh
stew ledger tail <ledger> [--limit N] [--path <dir>]
```

It prints the last N Stew entries from `.stew/<ledger>.md`, defaulting to 10 entries. Output contains entries only, in file order, so the selected slice remains chronological.

## Interface Changes
- Add `stew ledger tail <ledger> --limit 5`.
- Add help support for `stew ledger tail --help` and `stew help ledger tail`.
- Update `stew ledger --help` to list `tail        Print recent ledger entries`.
- `--limit` defaults to `10`; values must be positive integers.
- Output has no ledger H1 or managed marker, only selected entries.

## Implementation Changes
- Add `internal/stewledgertail` with `Run(Options) (Result, error)`.
- Reuse the existing known-ledger read path from `internal/stewledgercat.Run`, then parse entries from the raw markdown.
- Detect entries by canonical Stew H2 timestamp headings: lines like `## 2026-05-01T02:08:06Z — Summary`.
- Return the last N detected entries, preserving each selected entry’s markdown body and separators.
- If a known ledger has no entries, return success with empty stdout.
- Wire `tail` into `internal/cli/ledger.go` with `--limit` and `--path` flags.

## Test Plan
- Add package tests for:
  - returns the last N entries by H2 timestamp boundary
  - preserves file order for selected entries
  - omits the ledger title/managed marker
  - returns all entries when `--limit` exceeds entry count
  - returns empty content for a known ledger with no entries
  - rejects zero and negative limits
  - ignores non-entry H2 headings inside entry bodies
- Add CLI tests for:
  - `stew ledger tail iterations --limit 2`
  - default limit behavior
  - `stew ledger tail iterations --path <tmp repo>`
  - `stew ledger tail --help`
  - `stew help ledger tail`
  - `stew ledger --help` includes `tail`
- Verify with:
  - `go test ./...`
  - `go build ./...`
  - `go run ./cmd/stew ledger tail iterations --limit 5`

## Assumptions
- Default limit is `10`.
- Output is entries only.
- Output order is oldest-to-newest among the selected tail entries.
- Search and further filtering remain shell-pipe concerns.
</proposed_plan>

---

## 2026-05-01T02:38:47Z — Plan spec workflow move

**Prompt:** Plan moving Stew agent workflow into stew.spec.md and simplifying AGENTS.md.

<proposed_plan>
# Move Stew Workflow Into `stew.spec.md`

## Summary
Add a `## Working With Stew` section to `.stew/stew.spec.md` and the init template so agents get the full workflow from `stew full-spec`. Reduce the managed `AGENTS.md` block to a short pointer that says this repo uses Stew durable memory, then tells agents to run `stew help` and `stew full-spec`.

## Key Changes
- Update the Stew spec to define the agent workflow:
  - Start each new task/session with `stew help`, `stew full-spec`, and `stew ledgers`.
  - Tail every discovered writable ledger with `stew ledger tail <name> --limit 5`.
  - Use the tailed entries as recent project memory, while still verifying current repo state from files.
  - Before unfamiliar write commands, run command-specific help.
  - At the end of meaningful work, append to the appropriate ledgers according to their specs.
- Update the managed `AGENTS.md` block and `renderManagedBlock()` to keep only:
  - Stew manages durable project memory in append-only ledgers.
  - Run `stew help` for command discovery.
  - Run `stew full-spec` before working to load the workflow and ledger contract.
- Update the current `.stew/stew.spec.md`, current `AGENTS.md`, and README’s agent-first wording to match the new source-of-truth model.

## Test Plan
- Update init/template tests to assert:
  - generated AGENTS block is minimal and points to `stew help` + `stew full-spec`
  - generated `stew.spec.md` includes `Working With Stew`
  - workflow mentions `stew ledgers`, `stew ledger tail <name> --limit 5`, and end-of-work appends
  - AGENTS no longer lists default ledgers or detailed append/new-ledger instructions
- Verify with:
  - `go test ./...`
  - `go build ./...`
  - `go run ./cmd/stew full-spec`
  - `go run ./cmd/stew ledger tail iterations --limit 5`

## Assumptions
- Startup tailing should cover all discovered writable ledgers, not only defaults.
- Startup tail limit is 5 entries per ledger.
- `AGENTS.md` keeps both `stew help` and `stew full-spec`, but detailed workflow belongs only in `stew.spec.md`.
</proposed_plan>

---
