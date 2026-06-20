## 2026-06-20T22:39:14Z — Update Go toolchain pin

**Prompt:** PLEASE IMPLEMENT THIS PLAN: Update Stew To Go 1.26.4

Updated Stew's Go version from 1.22.0 to 1.26.4 in go.mod using go get go@1.26.4. Updated the GitHub release workflow to use actions/setup-go with go-version '1.26.4' so release artifacts build with the same exact supported toolchain.

Validation: go test ./...; go build ./...; make build; ./dist/stew version; ./dist/stew full-spec; rg -n "go 1\\.26\\.4|go-version: '1\\.26\\.4'" go.mod .github/workflows/release.yml; git diff --check.

---
