# Decisions

<!-- Managed by stew -->

## 2026-04-26T18:46:53Z — Canonicalize spec loading via full-spec

**Prompt:** Add to the existing init behavior by creating `.stew/stew.spec.md`, updating the AGENTS managed block to use `stew full-spec` as the canonical entry point, and making `stew full-spec` auto-include custom ledger specs.

**Context:** Agents previously depended on remembering multiple individual spec file paths, and extending with new ledgers required extra coordination in AGENTS guidance.

**Decision:** Treat `stew full-spec` as the single canonical contract loader. Add `.stew/stew.spec.md` as the base model spec, have `stew full-spec` concatenate all `.stew/*.spec.md` files (with `stew.spec.md` first), and point the AGENTS managed block to that single command.

**Consequences:** Session-start bootstrap for agents becomes one command, custom ledgers with `*.spec.md` are included automatically without AGENTS edits, and `stew init` now materializes six default files instead of five.

---

## 2026-04-26T23:15:36Z — Make ledger creation first-class

**Prompt:** PLEASE IMPLEMENT THIS PLAN: Implement `stew ledger new`.

**Context:** Custom ledgers were already discoverable through `.stew/*.spec.md`, but creating one still required hand-authoring both files and remembering the shape of a valid spec.

**Decision:** Treat `stew ledger new <name>` as the single first-class creation primitive for custom ledgers. The command creates `.stew/<name>.md` and `.stew/<name>.spec.md`, captures optional description and append-threshold guidance, validates canonical filesystem-safe names, and leaves discovery filesystem-based with no registry or config mutation.

**Consequences:** New ledgers become easy for agents and humans to scaffold consistently. List/remove/rename remain normal filesystem operations, existing repos need no migration, and strict creation-time name validation applies only to newly created ledgers while existing append compatibility remains unchanged.

---

## 2026-04-26T23:22:56Z — Create ledgers only in current repo

**Prompt:** i dont see a need to support custom file paths

**Context:** `stew ledger new` originally mirrored other commands by accepting `--path`, but ledger creation is most understandable when it works only in the current repository. A `--path` flag also makes the command look like it supports arbitrary custom ledger file locations, which is not the intended model.

**Decision:** Remove `--path` from the public `stew ledger new` command. Ledger creation now always resolves the current working directory, while the internal package keeps `TargetDir` for testability.

**Consequences:** The command surface is simpler and avoids implying custom file-path support. Scripts that need to create a ledger elsewhere should `cd` into that repository first. This intentionally diverges from `init`, `append`, and `full-spec`, which still keep `--path` for cross-repo operation.

---

## 2026-04-27T00:06:09Z — Store full AI plans in plans ledger

**Prompt:** awesome, shall we test it. I want to create a plans ledger that stores all the plans that I create with AI agents. I think it could be pretty cool to see all implemented plans after a while. I also think we should try to save it in entirety, whatever the model gave us plan to confirm, just save it. So we reduce work on model's end.

**Context:** Implementation plans produced by AI agents are useful historical artifacts, especially after they are implemented, but summarizing them again during logging wastes model effort and can lose intent.

**Decision:** Add a `plans` ledger for full AI-generated implementation plans. Entries should preserve the model-proposed plan in full as the body, including original structure and `<proposed_plan>` tags when present, with only an optional outcome note after the copied plan.

**Consequences:** Future sessions can review actual proposed plans rather than summaries. The ledger may grow quickly because full plans are stored verbatim, but that is intentional for auditability and lower logging friction.

---

## 2026-05-01T02:38:56Z — Use stew spec as workflow source

**Prompt:** Move Stew agent workflow into stew.spec.md and simplify AGENTS.md to a pointer.

**Context:** The managed AGENTS.md block had grown into detailed command guidance as new read and write commands were added. That made every initialized repo duplicate workflow instructions that belong in the Stew contract itself.

**Decision:** Treat `.stew/stew.spec.md` as the source of truth for the AI-agent Stew workflow. The spec now tells agents to run `stew help`, `stew full-spec`, `stew ledgers`, tail every discovered writable ledger with `stew ledger tail <name> --limit 5`, verify repo state from files, and append to the appropriate ledgers after meaningful work. The managed AGENTS.md block is reduced to a durable-memory pointer plus `stew help` and `stew full-spec` entrypoints.

**Consequences:** Future `stew init` output keeps AGENTS.md compact while still giving agents a complete workflow through `stew full-spec`. Workflow changes can be maintained in the Stew spec/template instead of duplicating command-level details in AGENTS.md. Agents are now explicitly expected to load recent context from all discovered ledgers at session start.

---

## 2026-05-01T02:43:26Z — Require append command writes

**Prompt:** Clarify that ledger paths from `stew ledgers` are not write targets.

**Context:** The Stew workflow now tells agents to run `stew ledgers`, whose output includes ledger file paths. Without an explicit write rule, those paths could be misread as permission to edit ledger markdown files directly.

**Decision:** Make the core model state that ledger file paths are for discovery and reading, and that entries must be created only with `stew append <ledger>`, never by directly editing `.stew/<ledger>.md`. The workflow section also repeats this before end-of-task writes.

**Consequences:** Agents can still use ledger paths for inspection, `cat`, `tail`, and shell tooling, but all ledger mutations stay behind the append command so Stew owns timestamps, prompt attribution, separators, and append-only formatting.

---

## 2026-05-01T02:51:34Z — Use CLI as ledger interface

**Prompt:** Treat Stew CLI commands as the only model-facing ledger interface; avoid exposing ledger file paths.

**Context:** The workflow originally said `stew ledgers` lists ledger paths and briefly described ledger paths as useful for discovery and reading. That exposed filesystem details to agents and created ambiguity about whether ledger files should be read or edited directly.

**Decision:** Treat Stew commands as the model-facing interface for ledger access. Agents should discover ledger names with `stew ledgers`, read with `stew ledger tail` or `stew ledger cat`, and write with `stew append`. The workflow and command help should avoid teaching agents to use ledger file paths, and `stew ledgers` should print names and descriptions rather than paths.

**Consequences:** Agents have a simpler operational model: use Stew for ledger operations, use repository files for source-code truth, and do not inspect or edit ledger files directly. Human-facing internals can still use the filesystem, but the agent contract no longer depends on path knowledge.

---
