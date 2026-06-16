## 2026-05-01T04:25:14Z — Clarify startup workflow spec

**Prompt:** Make Stew startup workflow explicit and remove redundant spec text

Made the Stew startup workflow more explicit by rewriting Working With Stew as a required context collection workflow after AGENTS.md points agents to full-spec. Removed the redundant Full Contract section and removed the iterations timestamp note because timestamps are produced by stew append. Updated the init template and tests so generated repos preserve the clearer contract. Validation: go test ./..., make build, and ./dist/stew full-spec | sed -n '1,90p'.

---
