# Stew Spec

Stew stores project context as append-only markdown ledgers managed through the `stew` CLI.
This file is the base model; ledger-specific requirements are included by `stew full-spec`.

## Core Model

A ledger has these durable properties:

- Name: the command-facing identifier, such as `iterations`.
- Spec: a ledger-specific contract defines the ledger's purpose, body conventions, and append threshold.
- Shared format: entries follow this base model plus the ledger-specific spec.
- Append-only semantics: never edit past entries.
- CLI interface: use `stew ledgers` to discover ledger names, `stew ledger tail` or `stew ledger cat` to read entries, and `stew append` to write entries. Do not inspect or edit ledger files directly.
- Chronological ordering: newest entries are appended at the bottom.
- Entry boundaries: each entry starts with an H2 UTC ISO 8601 timestamp and summary.
- Attribution: each entry includes `**Prompt:**` for the originating prompt.
- Repo affiliation: ledgers belong to the repository where Stew is initialized.

To add a custom ledger, run `stew ledger new <name>`. Stew commands automatically discover ledgers by their spec file.

## Full Contract

Run `stew full-spec` to print this file plus every ledger-specific spec.
Custom ledgers are discovered automatically when they are created with `stew ledger new`.

## Working With Stew

When starting a new task or session in a repository that uses Stew:

1. Run `stew help` to discover the available commands.
2. Run `stew full-spec` to load this workflow plus every ledger-specific contract.
3. Run `stew ledgers` to list writable ledger names and descriptions.
4. Run `stew ledger tail --all --limit 5` to load recent project memory from every discovered writable ledger.
5. Use the tailed entries as context, but verify current repo state from the actual files before changing behavior.

During the task, run `stew <command> --help` before using an unfamiliar write command.

At the end of meaningful work, append entries to the appropriate ledgers according to their specs. Use `iterations` for per-prompt work logs, `decisions` for durable architectural or product decisions, and any custom ledgers when their specs say the work belongs there.

Use Stew commands as the interface for ledger access: read with `stew ledger tail` or `stew ledger cat`, and write with `stew append`. Do not inspect or edit ledger files directly.
