## 2026-05-01T04:02:55Z — Remove synthetic tail IDs

**Prompt:** Remove synthetic tail entry IDs; IDs will become a primitive later

Removed the synthetic tail entry identifier from the internal entry model, JSON renderer, tests, and README example. Tail JSON now includes only timestamp, summary, prompt, and body for each entry; durable entry identifiers are left for a future markdown-level primitive. Validation: go test ./..., make build, and ./dist/stew ledger tail iterations --limit 1 --json with no id key.

---
