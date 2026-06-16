## 2026-04-26T23:22:56Z — Create ledgers only in current repo

**Prompt:** i dont see a need to support custom file paths

**Context:** `stew ledger new` originally mirrored other commands by accepting `--path`, but ledger creation is most understandable when it works only in the current repository. A `--path` flag also makes the command look like it supports arbitrary custom ledger file locations, which is not the intended model.

**Decision:** Remove `--path` from the public `stew ledger new` command. Ledger creation now always resolves the current working directory, while the internal package keeps `TargetDir` for testability.

**Consequences:** The command surface is simpler and avoids implying custom file-path support. Scripts that need to create a ledger elsewhere should `cd` into that repository first. This intentionally diverges from `init`, `append`, and `full-spec`, which still keep `--path` for cross-repo operation.

---
