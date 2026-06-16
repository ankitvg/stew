## 2026-05-01T03:59:39Z — Make tail JSON entry-aware

**Prompt:** This is trash. The JSON format we actually want is entry-aware

Replaced stew ledger tail --json content-string output with entry-aware JSON. Single-ledger JSON now returns ledger plus entries with id, timestamp, summary, prompt, and body fields; --all JSON returns sorted ledger objects with entry arrays. Added deterministic stew IDs, parser coverage for multiline prompts, strict CLI JSON shape tests, documentation, and rebuilt dist/stew. Validation: go test ./..., make build, ./dist/stew ledger tail iterations --limit 1 --json, and ./dist/stew ledger tail --all --limit 1 --json.

---
