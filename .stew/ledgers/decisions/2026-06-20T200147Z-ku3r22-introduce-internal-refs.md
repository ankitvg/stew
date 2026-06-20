## 2026-06-20T20:01:47Z — Introduce internal refs

**Prompt:** PLEASE IMPLEMENT THIS PLAN: Introduce Stew Refs As A Foundational Address Layer

**Context:** Stew now has stable atomic entry filenames and needs a small address vocabulary before adding relationship primitives such as links.

**Decision:** Add refs as an internal address layer with v1 support for `entry:<ledger>/<entry-file.md>` and `file:<repo-relative-path>`. Refs parse, canonicalize, and resolve project objects, but there is no public `stew ref` command yet.

**Consequences:** Future primitives can share a common address system. Append results now expose the exact generated entry ref internally, while normal CLI output remains unchanged.

---
