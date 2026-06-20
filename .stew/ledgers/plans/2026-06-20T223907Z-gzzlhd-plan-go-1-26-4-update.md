## 2026-06-20T22:39:07Z — Plan Go 1.26.4 update

**Prompt:** plan to go update to latest stable release

<proposed_plan>
# Update Stew To Go 1.26.4

## Summary
Update Stew from Go 1.22 to the current latest stable Go patch, `go1.26.4`, confirmed from the official [Go release history](https://go.dev/doc/devel/release). This is a toolchain maintenance change only; no Stew CLI behavior, public commands, schemas, or output contracts should change.

## Key Changes
- Update `go.mod` using the Go toolchain command:
  ```sh
  go get go@1.26.4
  ```
  Expected result:
  ```go
  go 1.26.4
  ```
- Update the GitHub release workflow to build/test/release with the same exact version:
  ```yaml
  go-version: '1.26.4'
  ```
- Do not add a `toolchain` directive unless the Go command unexpectedly writes one; a temp check showed `go get go@1.26.4` only updates the `go` directive.
- Do not cut a release, create a tag, or update Homebrew as part of this change.

## Test Plan
- Run:
  ```sh
  go test ./...
  go build ./...
  make build
  ```
- Confirm the local binary still works:
  ```sh
  ./dist/stew version
  ./dist/stew full-spec
  ```
- Confirm the release workflow still expresses the intended exact Go version:
  ```sh
  rg -n "go 1\\.26\\.4|go-version: '1\\.26\\.4'" go.mod .github/workflows/release.yml
  ```
- Run `git diff --check`.

## Assumptions
- Version policy is exact pinning: both source minimum and release CI use `1.26.4`.
- Go 1.26.4 is acceptable as the new minimum source-build requirement for Stew.
- No compatibility matrix is needed because Stew is a CLI application with no third-party Go module dependencies.
</proposed_plan>

**Outcome:** Implemented as planned. Updated `go.mod` and the GitHub release workflow to exact Go 1.26.4 pins, with no CLI behavior changes.

---
