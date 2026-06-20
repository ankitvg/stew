## 2026-06-20T19:17:22Z — Use ids for generated entry filenames

**Prompt:** yeah I think this is worth getting right now, add Stew-owned entropy to atomic entry filenames

**Context:** The initial entropy change covered `stew append`, but `stew migrate atomic-entries` and the legacy StewReads importer also create atomic entry files. Leaving those writers deterministic would make the storage contract uneven.

**Decision:** Every repo writer that creates new atomic entry files now uses compact timestamp, short Stew-generated lowercase base32 id, and readable summary slug in the filename. Existing historical deterministic filenames remain readable and valid.

**Consequences:** New append, migration, and importer writes reduce cross-clone filename collision risk while preserving chronological sort order. The older deterministic filename helper remains for compatibility-oriented tests and old entry parsing, not for new CLI writes.

---
