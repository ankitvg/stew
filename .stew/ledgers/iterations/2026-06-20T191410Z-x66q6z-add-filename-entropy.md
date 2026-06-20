## 2026-06-20T19:14:10Z — Add filename entropy

**Prompt:** yeah I think this is worth getting right now, add Stew-owned entropy to atomic entry filenames

Added short Stew-owned entry ids to new atomic append filenames. The append path now generates a six-character lowercase base32 id per new entry and writes filenames as compact timestamp, id, and slug, while retaining numeric suffix fallback for exact local collisions. Kept legacy deterministic filename helpers for migration and test fixtures.

Updated README and the Stew spec template/current spec to document the id-bearing filename contract. Added tests for id-bearing filename formatting, id validation, random id shape, append path output, and local suffix fallback.

Validation: go test ./internal/stewentry ./internal/stewappend; go test ./...; make build.

---
