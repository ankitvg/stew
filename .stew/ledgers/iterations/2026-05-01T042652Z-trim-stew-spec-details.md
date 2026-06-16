## 2026-05-01T04:26:52Z — Trim stew spec details

**Prompt:** Keep Stew spec lean; remove repeated CLI help details

Trimmed command-help-style details from the Stew spec and init template: removed the core-model CLI interface bullet, custom-ledger command paragraph, and final repeated CLI interface sentence. Kept only durable ledger properties plus the mandatory startup context workflow and append-at-end guidance. Updated tests to assert the leaner contract and forbid the removed details. Validation: go test ./..., make build, and ./dist/stew full-spec | sed -n '1,70p'.

---
