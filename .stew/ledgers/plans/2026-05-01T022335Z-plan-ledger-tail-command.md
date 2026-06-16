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
