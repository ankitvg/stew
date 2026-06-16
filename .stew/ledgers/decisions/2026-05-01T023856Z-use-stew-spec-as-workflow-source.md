## 2026-05-01T02:38:56Z — Use stew spec as workflow source

**Prompt:** Move Stew agent workflow into stew.spec.md and simplify AGENTS.md to a pointer.

**Context:** The managed AGENTS.md block had grown into detailed command guidance as new read and write commands were added. That made every initialized repo duplicate workflow instructions that belong in the Stew contract itself.

**Decision:** Treat `.stew/stew.spec.md` as the source of truth for the AI-agent Stew workflow. The spec now tells agents to run `stew help`, `stew full-spec`, `stew ledgers`, tail every discovered writable ledger with `stew ledger tail <name> --limit 5`, verify repo state from files, and append to the appropriate ledgers after meaningful work. The managed AGENTS.md block is reduced to a durable-memory pointer plus `stew help` and `stew full-spec` entrypoints.

**Consequences:** Future `stew init` output keeps AGENTS.md compact while still giving agents a complete workflow through `stew full-spec`. Workflow changes can be maintained in the Stew spec/template instead of duplicating command-level details in AGENTS.md. Agents are now explicitly expected to load recent context from all discovered ledgers at session start.

---
