## 2026-05-01T03:53:38Z — Add tail JSON output

**Prompt:** PLEASE IMPLEMENT THIS PLAN: Add --json To Ledger Tail

Added --json support to stew ledger tail with a unified ledgers array for single-ledger and --all output. The JSON renderer stays in the CLI layer, preserves tailed markdown content, omits filesystem paths, and documentation now includes JSON examples. Validated with go test ./... and make build.

---
