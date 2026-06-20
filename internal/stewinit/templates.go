package stewinit

const (
	managedBlockStart = "<!-- BEGIN STEW (managed) -->"
	managedBlockEnd   = "<!-- END STEW (managed) -->"
)

func renderManagedBlock() string {
	return managedBlockStart + "\n" +
		"## Stew\n\n" +
		"This repo uses stew to maintain durable project memory in append-only markdown ledger entries.\n\n" +
		"Run `stew help` to discover available commands.\n" +
		"Run `stew full-spec` before working to load the workflow and ledger contract.\n" +
		managedBlockEnd + "\n"
}

func renderStewSpec() string {
	return "# Stew Spec\n\n" +
		"Stew stores project context as named ledgers with append-only markdown entry files managed through the `stew` CLI.\n" +
		"This file is the base model; ledger-specific requirements are included by `stew full-spec`.\n\n" +
		"## Core Model\n\n" +
		"A ledger has these durable properties:\n\n" +
		"- Name: the command-facing identifier, such as `iterations`.\n" +
		"- Spec: a ledger-specific contract defines the ledger's purpose, body conventions, and append threshold.\n" +
		"- Storage: each entry is its own markdown file under `.stew/ledgers/<ledger>/`; specs stay at `.stew/<ledger>.spec.md`.\n" +
		"- Shared format: entries follow this base model plus the ledger-specific spec.\n" +
		"- Append-only semantics: never edit past entry files.\n" +
		"- Chronological ordering: entry filenames begin with a compact UTC timestamp and sort oldest-to-newest; newly generated filenames include a short Stew-generated id before the summary slug.\n" +
		"- Entry boundaries: each entry starts with an H2 UTC ISO 8601 timestamp and summary.\n" +
		"- Attribution: each entry includes `**Prompt:**` for the originating prompt.\n" +
		"- Repo affiliation: ledgers belong to the repository where Stew is initialized.\n\n" +
		"Older repositories with `.stew/<ledger>.md` files must be upgraded explicitly with `stew migrate atomic-entries`.\n\n" +
		"## Refs\n\n" +
		"Refs are canonical strings Stew uses to address project objects.\n\n" +
		"Stew currently supports two ref kinds:\n\n" +
		"- `entry:<ledger>/<entry-file.md>` addresses an atomic ledger entry.\n" +
		"- `file:<repo-relative-path>` addresses a repository file.\n\n" +
		"```txt\n" +
		"entry:decisions/2026-06-20T191722Z-5cxsdb-use-ids-for-generated-entry-filenames.md\n" +
		"file:internal/stewentry/stewentry.go\n" +
		"```\n\n" +
		"JSON tail output includes entry refs for returned entries.\n\n" +
		"## Links\n\n" +
		"Links are append-only relationships between two refs. A link has a source ref, a target ref, and a creation timestamp.\n\n" +
		"Stew stores links as JSON files under `.stew/links/`. V1 links have no kind field.\n\n" +
		"`stew append --link-file <repo-relative-path>` creates links from the new entry ref to file refs. `stew link list <ref>` lists links where the ref is the source or target.\n\n" +
		"## Working With Stew\n\n" +
		"Stew loads recent decisions and implementation notes so you can aim repo\n" +
		"inspection at the relevant files sooner instead of scanning blindly. Treat the\n" +
		"ledger context as a starting map, then confirm behavior against the current\n" +
		"source before making changes.\n\n" +
		"Required startup context workflow: after `AGENTS.md` tells you to run\n" +
		"`stew full-spec`, follow these steps before planning or editing for a new task\n" +
		"or session:\n\n" +
		"1. Run `stew ledgers` to list writable ledger names and descriptions.\n" +
		"2. Run `stew ledger tail --all --limit 5` to load recent project memory.\n" +
		"3. Use the tailed entries as context, then verify current repo state from the actual files before changing behavior.\n\n" +
		"During the task, run `stew <command> --help` before using an unfamiliar write command.\n\n" +
		"At the end of meaningful work, append entries to the appropriate ledgers according to their specs. Use `iterations` for per-prompt work logs, `decisions` for durable architectural or product decisions, and any custom ledgers when their specs say the work belongs there.\n"
}

func renderIterationsSpec() string {
	return `# Iterations Spec

A per-prompt work log. One entry per prompt that produced meaningful work.
Append-only.

## Entry shape

## <UTC ISO 8601, second precision> — <verb-first summary, 3-8 words>

**Prompt:** <the originating user prompt, verbatim or lightly compressed>

<Freeform body. Cover what was done, what changed, and how it was validated.
Write for a future reader reconstructing the work, not for a reviewer.>

---

## Notes

- One entry per prompt. If a prompt produced no meaningful change, no entry.
- Newest at bottom. Never edit past entries.
`
}

func renderDecisionsSpec() string {
	return `# Decisions Spec

Durable architectural or product decisions worth documenting. An entry belongs
here when the choice affects system architecture, product behavior, schema,
contracts, or anything a future engineer would need to understand to avoid
re-litigating it. Routine work goes in iterations, not here. Append-only.

## Entry shape

## <UTC ISO 8601, second precision> — <summary of the decision>

**Prompt:** <the originating user prompt or question>

**Context:** <what made this decision necessary>

**Decision:** <what was decided, stated plainly>

**Consequences:** <what changes, what's now ruled out>

---

## Notes

- Same timestamp and ordering rules as iterations.
- If a decision is later reversed, append a new entry that references the original. Never edit the original.
- When in doubt between iterations and decisions: if the next engineer needs to know this to avoid breaking something or repeating a debate, it's a decision.
`
}
