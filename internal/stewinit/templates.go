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
	return "# Iterations Spec\n\n" +
		"Use this format for every entry in `.stew/iterations.md`.\n\n" +
		"## Entry Format\n\n" +
		"1. Start each entry with an H2 timestamp in UTC ISO 8601 format with second precision.\n" +
		"2. Include a `**Prompt:**` line containing only the original user prompt.\n" +
		"3. Write a freeform markdown body with implementation details, outcomes, and follow-ups.\n" +
		"4. End each entry with a separator line: `---`.\n\n" +
		"## Example\n\n" +
		"## 2026-04-26T16:42:29Z\n\n" +
		"**Prompt:** Add OAuth retries to token exchange\n\n" +
		"Implemented retry-on-429 with capped exponential backoff and added integration coverage.\n\n" +
		"---\n"
}

func renderDecisionsSpec() string {
	return "# Decisions Spec\n\n" +
		"Use this format for every entry in `.stew/decisions.md`.\n\n" +
		"## Entry Format\n\n" +
		"1. Start each entry with an H2 timestamp in UTC ISO 8601 format with second precision.\n" +
		"2. Include a `**Prompt:**` line containing only the original user prompt.\n" +
		"3. Use decision sections: `**Context:**`, `**Decision:**`, and `**Consequences:**` (recommended convention).\n" +
		"4. End each entry with a separator line: `---`.\n\n" +
		"## Example\n\n" +
		"## 2026-04-26T16:42:29Z\n\n" +
		"**Prompt:** Choose default retry strategy for API client\n\n" +
		"**Context:** External provider latency spikes during peak windows.\n\n" +
		"**Decision:** Use bounded exponential backoff with jitter and a 15s overall timeout.\n\n" +
		"**Consequences:** Fewer transient failures; slightly higher average response time during incidents.\n\n" +
		"---\n"
}
