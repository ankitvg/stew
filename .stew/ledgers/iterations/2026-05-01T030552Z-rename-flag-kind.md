## 2026-05-01T03:05:52Z — Rename flag kind

**Prompt:** yeah argFlag sounds better, fix it

Renamed the internal CLI flag kind from `stringFlag` to `argFlag` so the parser reads as a flag that consumes an argument rather than a string-specific parser. Updated all parseInterspersedFlags call sites and preserved the existing Go flag typed parsing, including `IntVar` for ledger tail limits.

Validation:
- `go test ./...`
- `go build ./...`

---
