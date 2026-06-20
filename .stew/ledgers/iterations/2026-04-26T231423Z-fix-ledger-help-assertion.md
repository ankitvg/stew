## 2026-04-26T23:14:23Z — Fix ledger help assertion

**Prompt:** PLEASE IMPLEMENT THIS PLAN: Implement `stew ledger new`.

Initial `go test ./...` failed in `TestLedgerNewHelpDocumentsCreationContract` because Cobra wrapped the help text between `if` and `omitted`, making a single substring assertion too brittle. The implementation behavior was intact; the next step is to assert smaller stable fragments of the same help contract.

---
