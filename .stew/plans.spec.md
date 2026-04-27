# Plans Spec

Complete AI-agent implementation plans, preserved for later review after they have been proposed, confirmed, or executed.

## Body

Paste the full AI-generated plan as the entry body. Preserve the plan's original structure, headings, bullets, code blocks, and `<proposed_plan>` tags when present. Do not summarize, normalize, or rewrite the plan; the value of this ledger is keeping the proposed plan as the historical artifact.

If useful after implementation, add a brief `**Outcome:**` note after the copied plan. Never replace or compress the original plan text.

## Entry shape

The Stew header and prompt are produced by `stew append`:

```md
## <UTC ISO 8601, second precision> — <short plan summary>

**Prompt:** <prompt that produced the plan>

<full AI-generated plan exactly as proposed>
```

## When to append

Append whenever an AI agent produces an implementation plan for confirmation or execution; preserve the proposed plan in full instead of summarizing it.
