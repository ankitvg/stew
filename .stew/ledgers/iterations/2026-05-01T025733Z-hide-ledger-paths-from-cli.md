## 2026-05-01T02:57:33Z — Hide ledger paths from CLI

**Prompt:** Remove ledger paths from `stew ledgers` and ensure ledger paths are never shown through commands.

Updated the CLI and docs so ledger file paths are not part of the model-facing interface.

Changes:
- Confirmed `stew ledgers` prints only ledger names and descriptions.
- Changed `stew append` success output to `Appended <ledger>`.
- Changed `stew ledger new` success output to `Created ledger <name>`.
- Changed `stew init` output to summarize `Stew metadata` instead of listing managed files.
- Reworded help text and errors across read/write commands to avoid exposing ledger/spec file paths.
- Changed `stew full-spec` section comments to logical labels such as `<!-- stew spec -->` instead of filesystem paths.
- Updated README and tests to describe Stew commands as the interface for ledger access.

Validation:
- `go test ./...`
- `go build ./...`
- smoke test in a temp repo covering `init`, `ledger new`, `append`, `ledgers`, `full-spec`, and `ledger tail`, with output checked for no `.stew/` ledger paths.

---
