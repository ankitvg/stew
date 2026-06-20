package stewlink

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/ankitvg/stew/internal/stewentry"
	"github.com/ankitvg/stew/internal/stewref"
)

func TestCreateWritesCanonicalLink(t *testing.T) {
	tmp, entryRef := setupLinkRepo(t)
	targetRef, err := stewref.Parse(`file:./internal\cli/ledger.go`)
	if err != nil {
		t.Fatalf("parse target ref: %v", err)
	}

	result, err := Create(CreateOptions{
		TargetDir: tmp,
		Source:    entryRef,
		Target:    targetRef,
		Now:       fixedLinkNow,
		NewID:     staticLinkID("k7p3qx"),
	})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	if result.LinkPath != filepath.ToSlash(filepath.Join(".stew", "links", "2026-05-01T120507Z-k7p3qx.json")) {
		t.Fatalf("LinkPath = %q", result.LinkPath)
	}
	want := Link{
		Version:   Version,
		CreatedAt: "2026-05-01T12:05:07Z",
		Source:    entryRef.String(),
		Target:    "file:internal/cli/ledger.go",
	}
	if result.Link != want {
		t.Fatalf("Link = %#v, want %#v", result.Link, want)
	}

	stored := readStoredLink(t, filepath.Join(tmp, result.LinkPath))
	if stored != want {
		t.Fatalf("stored link = %#v, want %#v", stored, want)
	}
}

func TestCreateRejectsInvalidRefs(t *testing.T) {
	tmp, entryRef := setupLinkRepo(t)

	_, err := Create(CreateOptions{
		TargetDir: tmp,
		Source:    entryRef,
		Target:    stewref.Ref{Kind: "unknown", Payload: "value"},
		Now:       fixedLinkNow,
		NewID:     staticLinkID("k7p3qx"),
	})
	if !errors.Is(err, stewref.ErrInvalidRef) {
		t.Fatalf("Create() error = %v, want ErrInvalidRef", err)
	}
}

func TestCreateRequiresExistingSourceAndTarget(t *testing.T) {
	t.Run("source", func(t *testing.T) {
		tmp, _ := setupLinkRepo(t)
		missingSource, err := stewref.Entry("iterations", "2026-05-01T120507Z-k7p3qx-missing.md")
		if err != nil {
			t.Fatalf("entry ref: %v", err)
		}
		targetRef, err := stewref.File("internal/cli/ledger.go")
		if err != nil {
			t.Fatalf("file ref: %v", err)
		}

		_, err = Create(CreateOptions{
			TargetDir: tmp,
			Source:    missingSource,
			Target:    targetRef,
			Now:       fixedLinkNow,
			NewID:     staticLinkID("k7p3qx"),
		})
		if !errors.Is(err, stewref.ErrMissingTarget) {
			t.Fatalf("Create() error = %v, want ErrMissingTarget", err)
		}
	})

	t.Run("target", func(t *testing.T) {
		tmp, entryRef := setupLinkRepo(t)
		targetRef, err := stewref.File("internal/cli/missing.go")
		if err != nil {
			t.Fatalf("file ref: %v", err)
		}

		_, err = Create(CreateOptions{
			TargetDir: tmp,
			Source:    entryRef,
			Target:    targetRef,
			Now:       fixedLinkNow,
			NewID:     staticLinkID("k7p3qx"),
		})
		if !errors.Is(err, stewref.ErrMissingTarget) {
			t.Fatalf("Create() error = %v, want ErrMissingTarget", err)
		}
	})
}

func TestListReturnsLinksBySourceAndTarget(t *testing.T) {
	tmp, entryRef := setupLinkRepo(t)
	secondPath := filepath.Join(tmp, "internal", "stewlink", "stewlink.go")
	writeLinkTestFile(t, secondPath, "package stewlink\n")
	firstTarget, err := stewref.File("internal/cli/ledger.go")
	if err != nil {
		t.Fatalf("file ref: %v", err)
	}
	secondTarget, err := stewref.File("internal/stewlink/stewlink.go")
	if err != nil {
		t.Fatalf("file ref: %v", err)
	}

	first, err := Create(CreateOptions{
		TargetDir: tmp,
		Source:    entryRef,
		Target:    firstTarget,
		Now:       fixedLinkNow,
		NewID:     staticLinkID("aaaaaa"),
	})
	if err != nil {
		t.Fatalf("Create(first) error = %v", err)
	}
	second, err := Create(CreateOptions{
		TargetDir: tmp,
		Source:    entryRef,
		Target:    secondTarget,
		Now: func() time.Time {
			return fixedLinkNow().Add(time.Second)
		},
		NewID: staticLinkID("bbbbbb"),
	})
	if err != nil {
		t.Fatalf("Create(second) error = %v", err)
	}

	bySource, err := List(ListOptions{TargetDir: tmp, Ref: entryRef})
	if err != nil {
		t.Fatalf("List(source) error = %v", err)
	}
	if len(bySource.Links) != 2 {
		t.Fatalf("source links len = %d, want 2", len(bySource.Links))
	}
	if bySource.Links[0] != first.Link || bySource.Links[1] != second.Link {
		t.Fatalf("source links = %#v", bySource.Links)
	}

	byTarget, err := List(ListOptions{TargetDir: tmp, Ref: secondTarget})
	if err != nil {
		t.Fatalf("List(target) error = %v", err)
	}
	if len(byTarget.Links) != 1 || byTarget.Links[0] != second.Link {
		t.Fatalf("target links = %#v, want second", byTarget.Links)
	}
}

func TestListReturnsEmptyWhenNoLinksExist(t *testing.T) {
	tmp, _ := setupLinkRepo(t)
	ref, err := stewref.File("missing-but-syntactically-valid.go")
	if err != nil {
		t.Fatalf("file ref: %v", err)
	}

	result, err := List(ListOptions{TargetDir: tmp, Ref: ref})
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}
	if result.Ref != "file:missing-but-syntactically-valid.go" {
		t.Fatalf("Ref = %q", result.Ref)
	}
	if len(result.Links) != 0 {
		t.Fatalf("Links len = %d, want 0", len(result.Links))
	}
}

func TestListRejectsMalformedStoredLinkRecords(t *testing.T) {
	tmp, _ := setupLinkRepo(t)
	linksDir := filepath.Join(tmp, ".stew", "links")
	if err := os.MkdirAll(linksDir, 0o755); err != nil {
		t.Fatalf("mkdir links: %v", err)
	}
	writeLinkTestFile(t, filepath.Join(linksDir, "2026-05-01T120507Z-k7p3qx.json"), `{"version":0}`)
	ref, err := stewref.File("internal/cli/ledger.go")
	if err != nil {
		t.Fatalf("file ref: %v", err)
	}

	_, err = List(ListOptions{TargetDir: tmp, Ref: ref})
	if !errors.Is(err, ErrMalformedLink) {
		t.Fatalf("List() error = %v, want ErrMalformedLink", err)
	}
}

func setupLinkRepo(t *testing.T) (string, stewref.Ref) {
	t.Helper()
	tmp := t.TempDir()
	stewDir := filepath.Join(tmp, ".stew")
	if err := os.MkdirAll(filepath.Join(stewDir, "ledgers", "iterations"), 0o755); err != nil {
		t.Fatalf("mkdir ledger: %v", err)
	}
	if err := os.WriteFile(filepath.Join(stewDir, "iterations.spec.md"), []byte("# Iterations\n"), 0o644); err != nil {
		t.Fatalf("write spec: %v", err)
	}

	entryName := "2026-05-01T120507Z-k7p3qx-add-links.md"
	entryPath := filepath.Join(stewDir, "ledgers", "iterations", entryName)
	entry := stewentry.Render(fixedLinkNow(), "Add links", "Prompt", "Body")
	writeLinkTestFile(t, entryPath, entry)
	writeLinkTestFile(t, filepath.Join(tmp, "internal", "cli", "ledger.go"), "package cli\n")

	ref, err := stewref.Entry("iterations", entryName)
	if err != nil {
		t.Fatalf("entry ref: %v", err)
	}
	return tmp, ref
}

func fixedLinkNow() time.Time {
	return time.Date(2026, 5, 1, 12, 5, 7, 0, time.UTC)
}

func staticLinkID(id string) func() (string, error) {
	return func() (string, error) {
		return id, nil
	}
}

func readStoredLink(t *testing.T, path string) Link {
	t.Helper()
	bytes, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	var link Link
	if err := json.Unmarshal(bytes, &link); err != nil {
		t.Fatalf("unmarshal link: %v\n%s", err, string(bytes))
	}
	return link
}

func writeLinkTestFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("mkdir %s: %v", filepath.Dir(path), err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}
