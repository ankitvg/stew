# stew

Stew is a small CLI for maintaining append-only markdown ledger entries in a
repository. It is meant to give humans and coding agents a durable project
memory without a database, service, or generated registry.

Stew exposes project memory through named ledgers. Each ledger has entries and a
matching spec that explains when and how to write to it.

Ledger specs live at `.stew/<ledger>.spec.md`. Entries are stored atomically as
one markdown file per entry under `.stew/ledgers/<ledger>/`, with timestamped
filenames that sort oldest-to-newest. Newly generated filenames include short
Stew-generated ids before the summary slug.

Stew also has an internal ref vocabulary for addressing project objects. It
currently supports ledger entry refs, such as
`entry:decisions/2026-06-20T191722Z-5cxsdb-use-ids-for-generated-entry-filenames.md`,
and repo file refs, such as `file:internal/stewentry/stewentry.go`.

## Agent-First Usage

Stew is designed so most day-to-day usage is handled by an AI coding agent, not
by a human memorizing commands.

The setup loop is intentionally small: install the CLI, then run `stew init` in
a repo, or ask your coding agent to run it. Init creates Stew metadata and adds a
managed Stew block to `AGENTS.md`, so future agent sessions know to run
`stew help` and load `stew full-spec`. The full spec carries the agent workflow:
discover ledgers, tail recent entries from all ledgers for context, and append
to the appropriate ledgers after meaningful work.

After that, the user usually only needs to ask the agent to keep project context
up to date. The CLI remains documented and scriptable for anyone who wants to
drive it directly.

## Default Ledgers

`stew init` creates two ledgers by default:

- `iterations` is the per-prompt work log. Agents append here after meaningful
  work so future sessions can reconstruct what changed, why it changed, and how
  it was validated.
- `decisions` records durable architectural or product choices. Use it when a
  choice affects system behavior, contracts, schemas, or future tradeoffs that
  should not be re-litigated from scratch.

Routine implementation notes belong in `iterations`; decisions that future
maintainers need to preserve belong in `decisions`.

## Install

For a tagged release:

```sh
go install github.com/ankitvg/stew/cmd/stew@v0.1.0
```

From a local checkout:

```sh
go install ./cmd/stew
```

To build a local binary with version metadata:

```sh
make build VERSION=v0.1.0
./dist/stew version
```

## Quick Start

Initialize Stew in a git repository:

```sh
stew init
```

Load the full ledger contract before writing:

```sh
stew full-spec
```

List available ledgers:

```sh
stew ledgers
stew ledgers --json
```

Print a ledger for reading or shell filtering:

```sh
stew ledger cat iterations
stew ledger cat --all
stew ledger cat iterations | grep 'Prompt'
```

Print recent ledger entries:

```sh
stew ledger tail iterations --limit 5
stew ledger tail iterations --json --limit 5
stew ledger tail --all --limit 5
stew ledger tail --all --json --limit 5
```

JSON tail output is entry-aware:

```json
{
  "ledger": "iterations",
  "entries": [
    {
      "timestamp": "2026-05-01T03:53:38Z",
      "summary": "Add tail JSON output",
      "prompt": "PLEASE IMPLEMENT THIS PLAN: Add --json To Ledger Tail",
      "body": "Added --json support to stew ledger tail..."
    }
  ]
}
```

Append a work-log entry:

```sh
printf 'Implemented the parser change and ran go test ./...' \
  | stew append iterations \
      --prompt 'Fix parser edge case' \
      --summary 'Fix parser edge case'
```

Migrate an older repo from monolithic `.stew/<ledger>.md` files to atomic entry
files:

```sh
stew migrate atomic-entries --dry-run
stew migrate atomic-entries
```

Create a custom ledger:

```sh
stew ledger new plans \
  --description 'Reasoning artifacts for future work.' \
  --threshold 'Append when a plan captures durable intent or tradeoffs.'
```

Then append to it:

```sh
stew append plans \
  --prompt 'Plan release work' \
  --summary 'Record release plan' \
  -m 'Ship README, build metadata, and a v0.1.0 tag.'
```

## Commands

- `stew init` creates Stew metadata, default ledger storage/specs, and a managed `AGENTS.md` block.
- `stew help` prints the CLI workflow and available commands.
- `stew full-spec` prints the base Stew spec plus every custom ledger spec.
- `stew ledgers` lists discovered writable ledger names and descriptions; use `--json` for machine-readable output.
- `stew ledger cat <ledger>` prints one ledger's concatenated entry markdown; `--all` prints every ledger under name sections.
- `stew ledger tail <ledger>` prints recent entries from one ledger; `--all` prints recent entries from every ledger, and `--json` prints machine-readable output.
- `stew append <ledger>` appends a timestamped entry to a known ledger.
- `stew ledger new <name>` creates custom ledger storage and a spec.
- `stew migrate atomic-entries` splits legacy monolithic ledger files into atomic entry files.
- `stew version` prints build metadata.

## Release Check

Before tagging a release, run:

```sh
make pre-release VERSION=v0.1.0
./dist/stew version
```

`make pre-release` runs tests, verifies all packages build, and creates a
versioned local binary at `dist/stew`.

## License

Stew is licensed under the Apache License, Version 2.0. See [LICENSE](LICENSE).
