# Decisions Spec

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
