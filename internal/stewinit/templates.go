package stewinit

const (
	managedBlockStart = "<!-- BEGIN STEW (managed) -->"
	managedBlockEnd   = "<!-- END STEW (managed) -->"
)

func renderManagedBlock() string {
	return managedBlockStart + "\n" +
		"## Stew\n\n" +
		"This section is managed by `stew init`.\n\n" +
		"- Iterations ledger: `.stew/iterations.md`\n" +
		"- Decisions ledger: `.stew/decisions.md`\n" +
		"- Iterations entry spec: `.stew/iterations.spec.md`\n" +
		"- Decisions entry spec: `.stew/decisions.spec.md`\n\n" +
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
