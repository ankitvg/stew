# Stew Spec

Stew stores project context as append-only markdown ledgers under `.stew/`.
This file is the base model; ledger-specific requirements live in `*.spec.md` files.

## Core Model

A ledger has these durable properties:

- Name: the filename stem, such as `iterations` for `.stew/iterations.md`.
- Spec: `.stew/<name>.spec.md` defines the ledger's purpose, body conventions, and append threshold.
- Shared format: entries follow this base model plus the ledger-specific spec.
- Append-only semantics: never edit past entries.
- Chronological ordering: newest entries are appended at the bottom.
- Entry boundaries: each entry starts with an H2 UTC ISO 8601 timestamp and summary.
- Attribution: each entry includes `**Prompt:**` for the originating prompt.
- Repo affiliation: ledgers belong to the repository containing their `.stew/` directory.

To add a custom ledger, run `stew ledger new <name>`. Stew commands automatically discover ledgers by their spec file.

## Full Contract

Run `stew full-spec` to print this file plus every `.stew/*.spec.md` file.
Custom ledgers are discovered automatically when their spec files exist.
