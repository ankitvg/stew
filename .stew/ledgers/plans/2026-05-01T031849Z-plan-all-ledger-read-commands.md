## 2026-05-01T03:18:49Z — Plan all-ledger read commands

**Prompt:** PLEASE IMPLEMENT THIS PLAN: Add `--all` to `stew ledger cat` and `stew ledger tail`

# Add `--all` to `stew ledger cat` and `stew ledger tail`

## Summary
Add all-ledger read support:

```sh
stew ledger cat --all
stew ledger tail --all --limit 5
```

`--all` discovers writable ledgers with the same discovery path as `stew ledgers`, reads them in alphabetical ledger-name order, and prints each ledger under a command-facing name section. No ledger file paths are shown.

## Key Changes
- Update `stew ledger cat`:
  - Accept either `<ledger>` or `--all`, but not both.
  - Keep `stew ledger cat <ledger>` raw and unchanged.
  - Make `stew ledger cat --all` print all discovered ledgers as sections:

    ```md
    # decisions

    <raw decisions ledger content>

    # iterations

    <raw iterations ledger content>
    ```

- Update `stew ledger tail`:
  - Accept either `<ledger>` or `--all`, but not both.
  - Keep `stew ledger tail <ledger>` entry-only and unchanged.
  - Make `stew ledger tail --all --limit N` print the last N entries per discovered ledger, preserving chronological order within each ledger:

    ```md
    # decisions

    <last N decision entries>

    # iterations

    <last N iteration entries>
    ```

- Reuse `internal/stewledgers.List` for all-ledger discovery and existing `stewledgercat.Run` / `stewledgertail.Run` for reads.
- Add a small internal aggregation helper or package returning `{Name, Content}` sections so CLI formatting stays simple.
- Update help text for `stew ledger cat --help`, `stew ledger tail --help`, and `stew ledger --help`.
- Update `stew.spec.md`, init templates, and README workflow wording to prefer:

  ```sh
  stew ledger tail --all --limit 5
  ```

  for startup memory loading after `stew ledgers`.

## Test Plan
- Add CLI tests for:
  - `stew ledger cat --all`
  - `stew ledger tail --all --limit 2`
  - `--all --path <tmp repo>`
  - `--all` output sorted by ledger name
  - `--all` output uses ledger-name sections and never prints ledger paths
  - `<ledger> --all` returns a clear mutual-exclusion error
  - missing `<ledger>` without `--all` still errors
  - help output documents `--all`
- Add package/helper tests for:
  - all-ledger aggregation uses discovered ledger order
  - tail applies `--limit` per ledger, not globally
  - empty known ledgers produce an empty named section for `tail --all`
  - missing ledger content for any discovered ledger fails clearly
- Verify with:
  - `go test ./...`
  - `go build ./...`
  - `go run ./cmd/stew ledger cat --all`
  - `go run ./cmd/stew ledger tail --all --limit 5`

## Assumptions
- `--all` applies to both `cat` and `tail`.
- `--all` is mutually exclusive with a ledger positional.
- All-ledger output may add ledger-name section headings; single-ledger output remains unchanged.
- Ledger names are shown, but ledger file paths remain hidden from command output.

**Outcome:** Implemented as planned. Single-ledger `cat` and `tail` behavior remains unchanged, while `--all` renders ledger-name sections in discovered name order.

---
