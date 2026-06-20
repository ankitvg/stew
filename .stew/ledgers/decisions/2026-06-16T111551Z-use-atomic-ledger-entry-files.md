## 2026-06-16T11:15:51Z — Use atomic ledger entry files

**Prompt:** PLEASE IMPLEMENT THIS PLAN: Atomic Ledger Entry Storage

**Context:** Stew previously stored each logical ledger in one growing `.stew/<ledger>.md` file. That made writes simple, but large ledgers made storage, reviews, and future indexing heavier than necessary.

**Decision:** Store each ledger entry as its own markdown file under `.stew/ledgers/<ledger>/`, while keeping ledger specs flat at `.stew/<ledger>.spec.md` and preserving the existing command-facing ledger abstraction.

**Consequences:** `stew append` now writes one timestamped entry file, read commands concatenate or tail sorted entry files, and legacy monolithic ledgers require explicit `stew migrate atomic-entries`. Normal reads and appends no longer treat `.stew/<ledger>.md` as active storage.

---
