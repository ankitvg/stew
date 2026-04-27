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
