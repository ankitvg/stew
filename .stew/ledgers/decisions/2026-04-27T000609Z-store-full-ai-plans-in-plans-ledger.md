## 2026-04-27T00:06:09Z — Store full AI plans in plans ledger

**Prompt:** awesome, shall we test it. I want to create a plans ledger that stores all the plans that I create with AI agents. I think it could be pretty cool to see all implemented plans after a while. I also think we should try to save it in entirety, whatever the model gave us plan to confirm, just save it. So we reduce work on model's end.

**Context:** Implementation plans produced by AI agents are useful historical artifacts, especially after they are implemented, but summarizing them again during logging wastes model effort and can lose intent.

**Decision:** Add a `plans` ledger for full AI-generated implementation plans. Entries should preserve the model-proposed plan in full as the body, including original structure and `<proposed_plan>` tags when present, with only an optional outcome note after the copied plan.

**Consequences:** Future sessions can review actual proposed plans rather than summaries. The ledger may grow quickly because full plans are stored verbatim, but that is intentional for auditability and lower logging friction.

---
