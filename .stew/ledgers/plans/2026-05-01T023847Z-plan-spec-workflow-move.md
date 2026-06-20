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
