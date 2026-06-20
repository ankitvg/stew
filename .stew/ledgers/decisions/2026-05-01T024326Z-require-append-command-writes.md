## 2026-05-01T02:43:26Z — Require append command writes

**Prompt:** Clarify that ledger paths from `stew ledgers` are not write targets.

**Context:** The Stew workflow now tells agents to run `stew ledgers`, whose output includes ledger file paths. Without an explicit write rule, those paths could be misread as permission to edit ledger markdown files directly.

**Decision:** Make the core model state that ledger file paths are for discovery and reading, and that entries must be created only with `stew append <ledger>`, never by directly editing `.stew/<ledger>.md`. The workflow section also repeats this before end-of-task writes.

**Consequences:** Agents can still use ledger paths for inspection, `cat`, `tail`, and shell tooling, but all ledger mutations stay behind the append command so Stew owns timestamps, prompt attribution, separators, and append-only formatting.

---
