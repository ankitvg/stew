## 2026-06-20T18:47:57Z — Migrate StewReads repos to atomic ledgers

**Prompt:** i guess one off script makes sense, its just us using stewreads after all

Added a one-off StewReads legacy ledger importer at scripts/import_legacy_stewreads_ledgers.py. The script dry-runs by default, imports both monolithic .stew ledgers and pre-CLI docs ledgers, preserves historical timestamps where possible, normalizes legacy date-only/local-time headings, writes atomic entries under .stew/ledgers/<ledger>/, verifies counts and canonical headings, and removes legacy files only after successful verification when --remove-legacy is set.

Used it to migrate the four target StewReads repos on main: stew-docs (10 entries), stewreads-backend (200 entries), stewreads-ui (280 entries), and stewreads-mobile (22 entries). Each target repo was verified with ledgers, tail --all, and tail --all --json, then committed locally without pushing.

Validation passed: python3 -m py_compile scripts/import_legacy_stewreads_ledgers.py; dry-run across all four repos; make build; go test ./...; JSON tail checks for each migrated repo.

---
