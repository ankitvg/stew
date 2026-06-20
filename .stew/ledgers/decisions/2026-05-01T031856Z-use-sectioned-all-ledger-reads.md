## 2026-05-01T03:18:56Z — Use sectioned all-ledger reads

**Prompt:** Add all-ledger read support with --all

**Context:** `stew ledger cat` and `stew ledger tail` required a single ledger name, so agents had to discover ledgers and issue one read command per ledger to load project memory. The workflow needed a simpler command that still avoids exposing ledger file paths.

**Decision:** Add `--all` to `stew ledger cat` and `stew ledger tail`. `--all` is mutually exclusive with a ledger positional, discovers writable ledgers in name order, and renders each ledger under a command-facing `# <ledger>` section. Single-ledger output remains unchanged: `cat <ledger>` stays raw and `tail <ledger>` stays entry-only.

**Consequences:** Agents can load recent memory with `stew ledger tail --all --limit 5` after `stew ledgers`. Multi-ledger output is clearly attributable by ledger name without teaching file paths. Shell filtering still works on the combined stdout, while scripts depending on single-ledger raw output remain compatible.

---
