## 2026-04-28T03:17:58Z — Remove Cobra dependency

**Prompt:** PLEASE IMPLEMENT THIS PLAN: Remove Cobra And Use Standard Library CLI

Replaced the Cobra-based command layer with a standard-library dispatcher using flag.FlagSet, explicit help text, and ExecuteWithIO for tests. Preserved documented commands, root version/help flags, command-specific help, and append flags after the ledger positional. Removed Cobra, pflag, and mousetrap from go.mod and deleted go.sum because no third-party dependencies remain. Updated CLI tests away from Cobra command helpers and added a regression check that root help no longer lists completion. Validated with go test ./..., make pre-release VERSION=v0.1.0, manual help/version checks, a temp-repo append using post-positional flags, and an rg check for Cobra dependency residue.

---
