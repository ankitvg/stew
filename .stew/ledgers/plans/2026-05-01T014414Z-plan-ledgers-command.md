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
