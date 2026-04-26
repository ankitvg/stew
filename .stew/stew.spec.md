# Stew Spec

Stew stores project context as append-only markdown ledgers under `.stew/`.
This file is the base model; ledger-specific requirements live in `*.spec.md` files.

## Core Model

- Ledgers are markdown files in `.stew/` (for example `iterations.md`, `decisions.md`).
- Each ledger should have a companion spec file named `<ledger>.spec.md`.
- Entries are append-only. Never edit past entries.
- Use UTC ISO 8601 timestamps with second precision: `YYYY-MM-DDTHH:MM:SSZ`.

## Full Contract

Run `stew full-spec` to print this file plus every `.stew/*.spec.md` file.
Custom ledgers are discovered automatically when their spec files exist.
