## 2026-06-20T19:14:10Z — Use Stew generated filename ids

**Prompt:** yeah I think this is worth getting right now, add Stew-owned entropy to atomic entry filenames

**Context:** Atomic entry filenames were deterministic from timestamp plus summary. Local `O_EXCL` retries prevented same-checkout overwrites, but independent clones could still create the same path and hit git add/add conflicts.

**Decision:** New `stew append` entries include a short Stew-generated lowercase base32 id after the compact UTC timestamp and before the readable summary slug. Historical and migration filenames remain readable and deterministic.

**Consequences:** New append filenames preserve chronological sorting while reducing cross-clone filename collision risk. Numeric suffix retries remain as a local last-resort guard if an exact path already exists.

---
