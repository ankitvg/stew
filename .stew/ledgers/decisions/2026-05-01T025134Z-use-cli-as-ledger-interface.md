## 2026-05-01T02:51:34Z — Use CLI as ledger interface

**Prompt:** Treat Stew CLI commands as the only model-facing ledger interface; avoid exposing ledger file paths.

**Context:** The workflow originally said `stew ledgers` lists ledger paths and briefly described ledger paths as useful for discovery and reading. That exposed filesystem details to agents and created ambiguity about whether ledger files should be read or edited directly.

**Decision:** Treat Stew commands as the model-facing interface for ledger access. Agents should discover ledger names with `stew ledgers`, read with `stew ledger tail` or `stew ledger cat`, and write with `stew append`. The workflow and command help should avoid teaching agents to use ledger file paths, and `stew ledgers` should print names and descriptions rather than paths.

**Consequences:** Agents have a simpler operational model: use Stew for ledger operations, use repository files for source-code truth, and do not inspect or edit ledger files directly. Human-facing internals can still use the filesystem, but the agent contract no longer depends on path knowledge.

---
