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
