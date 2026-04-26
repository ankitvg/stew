package stewinit

const (
	managedBlockStart = "<!-- BEGIN STEW (managed) -->"
	managedBlockEnd   = "<!-- END STEW (managed) -->"
)

func renderManagedBlock() string {
	return managedBlockStart + "\n" +
		"## Stew\n\n" +
		"This repo uses stew to maintain append-only markdown ledgers.\n\n" +
		"Run `stew help` first to discover the CLI workflow and available commands.\n" +
		"Run `stew full-spec` to load the full contract (stew model + all ledger specs).\n" +
		"When a new ledger is needed, run `stew ledger new <name> --help`; then create it with `stew ledger new <name> ...`.\n" +
		"Before writing, run `stew append <ledger> --help`; then append entries with `stew append <ledger> ...`.\n\n" +
		"Default ledgers in this repo:\n" +
		"- iterations — per-prompt work log\n" +
		"- decisions — durable architectural and product decisions\n" +
		managedBlockEnd + "\n"
}

func renderIterationsLedger() string {
	return "# Iterations\n\n" +
		"<!-- Managed by stew -->\n"
}

func renderDecisionsLedger() string {
	return "# Decisions\n\n" +
		"<!-- Managed by stew -->\n"
}

func renderStewSpec() string {
	return "# Stew Spec\n\n" +
		"Stew stores project context as append-only markdown ledgers under `.stew/`.\n" +
		"This file is the base model; ledger-specific requirements live in `*.spec.md` files.\n\n" +
		"## Core Model\n\n" +
		"A ledger has these durable properties:\n\n" +
		"- Name: the filename stem, such as `iterations` for `.stew/iterations.md`.\n" +
		"- Spec: `.stew/<name>.spec.md` defines the ledger's purpose, body conventions, and append threshold.\n" +
		"- Shared format: entries follow this base model plus the ledger-specific spec.\n" +
		"- Append-only semantics: never edit past entries.\n" +
		"- Chronological ordering: newest entries are appended at the bottom.\n" +
		"- Entry boundaries: each entry starts with an H2 UTC ISO 8601 timestamp and summary.\n" +
		"- Attribution: each entry includes `**Prompt:**` for the originating prompt.\n" +
		"- Repo affiliation: ledgers belong to the repository containing their `.stew/` directory.\n\n" +
		"To add a custom ledger, run `stew ledger new <name>`. Stew commands automatically discover ledgers by their spec file.\n\n" +
		"## Full Contract\n\n" +
		"Run `stew full-spec` to print this file plus every `.stew/*.spec.md` file.\n" +
		"Custom ledgers are discovered automatically when their spec files exist.\n"
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

- Timestamp: ` + "`2026-04-26T18:32:00Z`" + `. No brackets, no abbreviations.
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
