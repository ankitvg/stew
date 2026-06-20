## 2026-06-20T22:11:54Z — Implement V1 links

**Prompt:** PLEASE IMPLEMENT THIS PLAN: Add V1 Links Without Kinds

Implemented V1 links as append-only source->target relationships between existing refs, without kind, link refs, delete/update, or tail embedding. Added internal link storage/listing under .stew/links, the read-only stew link list command, append --link-file creation with JSON output, init links storage creation, README/spec documentation, and focused tests.

Validation: go test ./internal/stewlink ./internal/stewappend ./internal/cli ./internal/stewinit; go test ./...; git diff --check; make build; manual temp-repo ./dist/stew init, append --link-file --json, and link list checks.

---
