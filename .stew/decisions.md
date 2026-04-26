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
