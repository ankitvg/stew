## 2026-05-01T03:30:46Z — Rebuild dist binary

**Prompt:** ankitgupta@Ankits-MacBook-Pro dist % ./stew ledgers --json\nunknown command "ledgers"

Rebuilt the local dist/stew binary after the JSON ledgers implementation because the existing dist binary was stale and did not include the ledgers command. Validation: make build, ./dist/stew ledgers --json, and ./dist/stew version.

---
