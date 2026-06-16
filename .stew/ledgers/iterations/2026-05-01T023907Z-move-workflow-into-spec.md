## 2026-05-01T02:39:07Z — Move workflow into spec

**Prompt:** PLEASE IMPLEMENT THIS PLAN: Move Stew Workflow Into `stew.spec.md`.

Moved the AI-agent workflow guidance into the Stew spec and reduced AGENTS.md to a compact pointer.

Changes:
- Added `## Working With Stew` to `.stew/stew.spec.md` and `renderStewSpec()` with startup instructions to run `stew help`, `stew full-spec`, `stew ledgers`, and `stew ledger tail <name> --limit 5` for each writable ledger.
- Updated the workflow to tell agents to use tailed entries as project memory while verifying current repo state from files, run command help before unfamiliar write commands, and append to appropriate ledgers after meaningful work.
- Simplified `AGENTS.md` and `renderManagedBlock()` so the managed block only describes Stew durable memory and points to `stew help` plus `stew full-spec`.
- Updated README agent-first wording to say the full spec carries the agent workflow.
- Updated init/template tests to enforce the minimal AGENTS block and spec workflow content.

Validation:
- `go test ./...`
- `go build ./...`
- `go run ./cmd/stew full-spec`
- `go run ./cmd/stew ledger tail iterations --limit 5`

---
