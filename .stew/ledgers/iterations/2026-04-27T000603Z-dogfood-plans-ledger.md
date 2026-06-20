## 2026-04-27T00:06:03Z — Dogfood plans ledger

**Prompt:** awesome, shall we test it. I want to create a plans ledger that stores all the plans that I create with AI agents. I think it could be pretty cool to see all implemented plans after a while. I also think we should try to save it in entirety, whatever the model gave us plan to confirm, just save it. So we reduce work on model's end.

Created a real `.stew/plans.md` and `.stew/plans.spec.md` in the stew repo using `stew ledger new plans` without `--path`. Tuned the generated plans spec so entries preserve the full AI-generated plan exactly as proposed, including headings, bullets, code blocks, and `<proposed_plan>` tags, with only an optional outcome note after the copied plan.

Appended the first real plans entry: the full previously confirmed `stew ledger new` implementation plan. Verified `stew full-spec` discovers the new plans spec and `go test ./...` still passes.

---
