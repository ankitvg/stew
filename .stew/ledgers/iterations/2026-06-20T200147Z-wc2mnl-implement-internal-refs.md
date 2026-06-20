## 2026-06-20T20:01:47Z — Implement internal refs

**Prompt:** PLEASE IMPLEMENT THIS PLAN: Introduce Stew Refs As A Foundational Address Layer

Implemented the internal Stew refs foundation in four logical commits: ref parsing/canonicalization, repo-aware resolution, append result entry refs, and refs documentation in README plus the generated/current spec.

The new `internal/stewref` package supports `entry:` and `file:` refs, validates unsafe paths and malformed entry refs, resolves refs under the repository root, and rejects `file:` refs into `.stew/ledgers/`. `stewappend.Result` now includes `EntryRef` built from the exact generated entry filename.

Validation: go test ./internal/stewref; go test ./internal/stewappend ./internal/stewinit ./internal/fullspec; go test ./...; make build; git diff --check; ./dist/stew full-spec inspection.

---
