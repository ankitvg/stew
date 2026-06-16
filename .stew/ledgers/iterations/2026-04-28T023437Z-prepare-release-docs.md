## 2026-04-28T02:34:37Z — Prepare release docs

**Prompt:** Lets skip license for a bit more? Do they rest

Expanded README with install, quick-start, command, and release-check guidance for v0.1.0. Updated Makefile so make build creates a versioned dist/stew binary with linker-injected version, commit, and date metadata, while make pre-release still runs tests and package builds first. Validated with make pre-release VERSION=v0.1.0 and ./dist/stew version.

---
