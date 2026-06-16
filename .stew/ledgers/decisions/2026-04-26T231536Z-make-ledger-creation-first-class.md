## 2026-04-26T23:15:36Z — Make ledger creation first-class

**Prompt:** PLEASE IMPLEMENT THIS PLAN: Implement `stew ledger new`.

**Context:** Custom ledgers were already discoverable through `.stew/*.spec.md`, but creating one still required hand-authoring both files and remembering the shape of a valid spec.

**Decision:** Treat `stew ledger new <name>` as the single first-class creation primitive for custom ledgers. The command creates `.stew/<name>.md` and `.stew/<name>.spec.md`, captures optional description and append-threshold guidance, validates canonical filesystem-safe names, and leaves discovery filesystem-based with no registry or config mutation.

**Consequences:** New ledgers become easy for agents and humans to scaffold consistently. List/remove/rename remain normal filesystem operations, existing repos need no migration, and strict creation-time name validation applies only to newly created ledgers while existing append compatibility remains unchanged.

---
