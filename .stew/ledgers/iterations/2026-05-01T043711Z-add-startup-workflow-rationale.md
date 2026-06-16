## 2026-05-01T04:37:11Z — Add startup workflow rationale

**Prompt:** Explain why the Stew startup workflow exists before listing steps

Added rationale to Working With Stew explaining that startup ledger reads provide recent decisions and implementation notes to guide targeted repo inspection before editing. Mirrored the wording in the init template and updated tests to preserve it. Validation: go test ./..., make build, and ./dist/stew full-spec | sed -n '1,55p'.

---
