## 2026-04-26T18:46:53Z — Canonicalize spec loading via full-spec

**Prompt:** Add to the existing init behavior by creating `.stew/stew.spec.md`, updating the AGENTS managed block to use `stew full-spec` as the canonical entry point, and making `stew full-spec` auto-include custom ledger specs.

**Context:** Agents previously depended on remembering multiple individual spec file paths, and extending with new ledgers required extra coordination in AGENTS guidance.

**Decision:** Treat `stew full-spec` as the single canonical contract loader. Add `.stew/stew.spec.md` as the base model spec, have `stew full-spec` concatenate all `.stew/*.spec.md` files (with `stew.spec.md` first), and point the AGENTS managed block to that single command.

**Consequences:** Session-start bootstrap for agents becomes one command, custom ledgers with `*.spec.md` are included automatically without AGENTS edits, and `stew init` now materializes six default files instead of five.

---
