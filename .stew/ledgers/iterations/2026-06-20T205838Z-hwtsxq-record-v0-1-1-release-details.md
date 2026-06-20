## 2026-06-20T20:58:38Z — Record v0.1.1 release details

**Prompt:** [$stew-release] record complete v0.1.1 release details

Released Stew v0.1.1 from main at d5fc9c4: https://github.com/ankitvg/stew/releases/tag/v0.1.1. The release workflow was fixed before tagging so artifacts are explicitly cross-compiled for macOS/Linux arm64/amd64. Downloaded assets were verified by file type and checksum. Updated ankitvg/homebrew-tap at 38707ad to serve v0.1.1 from verified release tarballs instead of source builds. Homebrew validation passed: brew style, brew audit, brew reinstall, brew test, and /opt/homebrew/bin/stew version. Removed repo-local dist/stew PATH overrides from ~/.zprofile and ~/.zshrc so clean login shells resolve /opt/homebrew/bin/stew.

---
