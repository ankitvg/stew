## 2026-04-26T20:22:34Z — Dogfood ledgers and document primitive

**Prompt:** Dogfood on stewreads-docs. Migrate existing iterations.md and decisions.md into .stew/ via stew append. Also add one sentence to stew.spec.md explicitly naming custom ledgers as a primitive.

Dogfooded Stew on `stew-docs` by building a local `/tmp/stew-dogfood` binary, running `stew init`, and migrating the existing real legacy entries from `docs/iterations.md` and `docs/decisions.md` into `.stew/iterations.md` and `.stew/decisions.md` through `stew append` stdin bodies. This validated the append format against real project entries rather than synthetic examples.

Made the custom-ledger primitive explicit by adding the sentence to the Stew spec template and current `.stew/stew.spec.md`: to add a custom ledger, create `.stew/<name>.md` and `.stew/<name>.spec.md`; Stew commands discover ledgers by their spec file. Added a test to lock that generated spec guidance.

---
