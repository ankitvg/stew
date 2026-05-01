# Stew Spec

Stew stores project context as append-only markdown ledgers managed through the `stew` CLI.
This file is the base model; ledger-specific requirements are included by `stew full-spec`.

## Core Model

A ledger has these durable properties:

- Name: the command-facing identifier, such as `iterations`.
- Spec: a ledger-specific contract defines the ledger's purpose, body conventions, and append threshold.
- Shared format: entries follow this base model plus the ledger-specific spec.
- Append-only semantics: never edit past entries.
- Chronological ordering: newest entries are appended at the bottom.
- Entry boundaries: each entry starts with an H2 UTC ISO 8601 timestamp and summary.
- Attribution: each entry includes `**Prompt:**` for the originating prompt.
- Repo affiliation: ledgers belong to the repository where Stew is initialized.

## Working With Stew

Required startup context workflow: after `AGENTS.md` tells you to run
`stew full-spec`, follow these steps before planning or editing for a new task
or session:

This loads recent decisions and implementation notes so you can aim repo
inspection at the relevant files sooner instead of scanning blindly. Treat the
ledger context as a starting map, then confirm behavior against the current
source before making changes.

1. Run `stew ledgers` to list writable ledger names and descriptions.
2. Run `stew ledger tail --all --limit 5` to load recent project memory.
3. Use the tailed entries as context, then verify current repo state from the actual files before changing behavior.

During the task, run `stew <command> --help` before using an unfamiliar write command.

At the end of meaningful work, append entries to the appropriate ledgers according to their specs. Use `iterations` for per-prompt work logs, `decisions` for durable architectural or product decisions, and any custom ledgers when their specs say the work belongs there.
