## 2026-06-16T11:15:28Z — Plan atomic ledger storage

**Prompt:** PLEASE IMPLEMENT THIS PLAN: Atomic Ledger Entry Storage

# Atomic Ledger Entry Storage

## Summary
Convert Stew from one markdown file per logical ledger to one markdown file per ledger entry, while preserving the existing command-facing ledger abstraction.

Use this storage layout:

```txt
.stew/
  config.toml
  stew.spec.md
  iterations.spec.md
  decisions.spec.md
  ledgers/
    iterations/
      2026-06-16T173012Z-add-atomic-entry-storage.md
    decisions/
      2026-06-16T173100Z-use-atomic-ledger-entries.md
```

Ledger specs remain at `.stew/<ledger>.spec.md`; entry files live under `.stew/ledgers/<ledger>/`.

## Key Changes
- Update `stew append <ledger>` to create one new entry file instead of appending to `.stew/<ledger>.md`.
- Keep the entry markdown body format unchanged: H2 timestamp summary, `**Prompt:**`, body, and `---`.
- Generate filenames as `<UTC timestamp compact>Z-<summary-slug>.md`, with numeric suffixes for same-second collisions.
- Update `stew ledger cat`, `stew ledger tail`, and `--json` to read entry files from `.stew/ledgers/<ledger>/` in chronological filename order.
- Update `stew ledger new <name>` and `stew init` to create ledger entry directories instead of monolithic ledger markdown files.
- Add explicit migration command:

```sh
stew migrate atomic-entries [--path <dir>] [--dry-run]
```

It parses existing `.stew/<ledger>.md` files, writes atomic entry files, verifies counts, then removes the old monolithic `.md` files unless `--dry-run` is set.

## Interface And Compatibility
- `stew ledgers` continues discovering ledgers from `.stew/*.spec.md`; no registry file is added.
- `.stew/<ledger>.md` is no longer an active storage format after migration.
- Migration is hard and explicit: normal `append`, `tail`, `cat`, and `init` do not silently upgrade old repos.
- If a command sees a spec but no `.stew/ledgers/<ledger>/` directory, return a clear error telling the user to run `stew migrate atomic-entries` or `stew init` for a fresh repo.
- Update `stew full-spec`, generated init templates, README, and CLI help to describe atomic entry storage instead of large ledger files.

## Test Plan
- Add storage/parser tests for reading atomic entries in chronological order, preserving entry markdown, and parsing prompts/bodies into JSON.
- Add append tests for filename generation, slug cleanup, same-second collision suffixes, missing ledger errors, and directory creation behavior.
- Add migration tests covering multi-entry ledgers, empty ledgers, custom ledgers, malformed/no-entry files, dry-run behavior, deletion of old `.stew/<ledger>.md`, and rollback-safe failure before deletion.
- Add CLI tests for `append`, `ledger cat`, `ledger tail`, `tail --json`, `tail --all --json`, `ledger new`, `init`, and `migrate atomic-entries`.
- Validate with `go test ./...` and `make build`.

## Assumptions
- Hard migration means old monolithic ledger files are removed after successful conversion.
- Specs remain flat at `.stew/<ledger>.spec.md` to preserve existing ledger discovery.
- Entry filenames prioritize human inspection and stable sorting over opaque IDs.
- The current repo should be migrated as part of implementation so Stew dogfoods the new format immediately.

**Outcome:** Implemented. Stew now writes and reads atomic entry files, includes explicit `stew migrate atomic-entries`, and this repo's own ledgers were migrated into `.stew/ledgers/<ledger>/`.

---
