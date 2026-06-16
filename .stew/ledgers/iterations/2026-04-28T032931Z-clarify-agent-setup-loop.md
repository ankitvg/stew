## 2026-04-28T03:29:31Z — Clarify agent setup loop

**Prompt:** nice but may be we should mention the entire setup. After all, stew init indeed plugs stew into AGENTS.md, so in a way all user has to do is init stew which cli agent can even do

Refined README.md's Agent-First Usage section to describe the full bootstrap flow: install the CLI, run or ask an agent to run stew init, and let the managed AGENTS.md block teach future agent sessions to run stew help, load full-spec, append entries, and create ledgers. Validated with go test ./....

---
