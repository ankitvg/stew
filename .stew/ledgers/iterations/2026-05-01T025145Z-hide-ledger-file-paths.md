## 2026-05-01T02:51:45Z — Hide ledger file paths

**Prompt:** Revise the workflow/spec so the model uses Stew commands instead of ledger file paths.

Revised the workflow and CLI surfaces so agents use Stew commands as the ledger interface instead of reasoning about ledger file paths.

Changes:
- Reworded `.stew/stew.spec.md` and `renderStewSpec()` so the core model describes command-facing ledger names and the CLI interface: `stew ledgers`, `stew ledger tail`, `stew ledger cat`, and `stew append`.
- Updated the workflow to say `stew ledgers` lists ledger names and descriptions, not paths.
- Changed `stew ledgers` output to omit ledger paths and print only names plus descriptions.
- Changed `stew full-spec` source markers from path comments to logical labels such as `<!-- stew spec -->` and `<!-- decisions spec -->`.
- Cleaned path-oriented wording from relevant command help and README command summaries.
- Updated tests for the new output and wording.

Validation:
- `go test ./...`
- `go build ./...`
- `go run ./cmd/stew help | sed -n '1,35p'`
- `go run ./cmd/stew full-spec | sed -n '1,50p'`
- `go run ./cmd/stew ledgers`

---
