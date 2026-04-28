# stew

Stew is a small CLI for maintaining append-only markdown ledgers in a repository.
It is meant to give humans and coding agents a durable project memory without a
database, service, or generated registry.

Stew stores ledgers under `.stew/`. Each ledger has a markdown file for entries
and a matching `.spec.md` file that explains when and how to write to it.

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

Append a work-log entry:

```sh
printf 'Implemented the parser change and ran go test ./...' \
  | stew append iterations \
      --prompt 'Fix parser edge case' \
      --summary 'Fix parser edge case'
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

- `stew init` creates `.stew/`, default ledgers, and a managed `AGENTS.md` block.
- `stew help` prints the CLI workflow and available commands.
- `stew full-spec` prints `.stew/stew.spec.md` plus every custom ledger spec.
- `stew append <ledger>` appends a timestamped entry to a known ledger.
- `stew ledger new <name>` creates a custom ledger and spec.
- `stew version` prints build metadata.

## Release Check

Before tagging a release, run:

```sh
make pre-release VERSION=v0.1.0
./dist/stew version
```

`make pre-release` runs tests, verifies all packages build, and creates a
versioned local binary at `dist/stew`.
