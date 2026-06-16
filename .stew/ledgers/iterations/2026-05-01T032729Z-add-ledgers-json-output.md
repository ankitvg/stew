## 2026-05-01T03:27:29Z — Add ledgers JSON output

**Prompt:** PLEASE IMPLEMENT THIS PLAN: Add JSON Output To stew ledgers

Implemented a boolean --json flag for stew ledgers. The JSON response is a top-level object with ledgers entries containing only name and description, preserving sorted discovery order and omitting ledger or target paths. Plain-text output remains unchanged. Updated command help and README docs, and added CLI tests for valid JSON, path flag support, sorted order, newline termination, and path-field exclusion. Validation: go test ./..., go build ./..., go run ./cmd/stew ledgers --json, go run ./cmd/stew ledgers --help, and go run ./cmd/stew ledgers --json --path /Users/ankitgupta/Documents/stewreads/stew.

---
