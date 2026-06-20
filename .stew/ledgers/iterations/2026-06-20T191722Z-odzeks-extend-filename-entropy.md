## 2026-06-20T19:17:22Z — Extend filename entropy

**Prompt:** yeah I think this is worth getting right now, add Stew-owned entropy to atomic entry filenames

Extended the filename entropy work beyond append. `stew migrate atomic-entries` now injects the same six-character lowercase base32 entry id when writing migrated atomic files, and the one-off StewReads legacy importer now uses the same timestamp-id-slug shape. Updated migration tests and CLI migration coverage to assert id-bearing filenames while keeping deterministic filename helpers for legacy-compatible fixtures.

Final validation: go test ./internal/stewentry ./internal/stewappend ./internal/stewmigrateatomic ./internal/cli; python3 -m py_compile scripts/import_legacy_stewreads_ledgers.py; go test ./...; make build; git diff --check.

---
