<!-- BEGIN STEW (managed) -->
## Stew

This repo uses stew to maintain append-only markdown ledgers.

Run `stew help` first to discover the CLI workflow and available commands.
Run `stew full-spec` to load the full contract (stew model + all ledger specs).
Before writing, run `stew append <ledger> --help`; then append entries with `stew append <ledger> ...`.

Default ledgers in this repo:
- iterations — per-prompt work log
- decisions — durable architectural and product decisions
<!-- END STEW (managed) -->
