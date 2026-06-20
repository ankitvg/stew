## 2026-06-20T21:17:27Z — Plumb refs into surfaces

**Prompt:** Can you plumb refs into existing surfaces first; do not design the link schema yet.

Carried atomic entry file provenance through ledger reads so tail results include canonical entry refs. Exposed refs in ledger tail JSON and added append --json for the entry ref created by a write. Updated README, CLI help, and the generated/current Stew spec.

Validation: go test ./internal/stewledgercat ./internal/stewledgertail ./internal/stewledgerall ./internal/cli ./internal/stewinit ./internal/fullspec; go test ./...; git diff --check; make build; manual ./dist/stew full-spec, ledger tail --json, and append --json checks.

---
