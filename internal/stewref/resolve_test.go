package stewref

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestResolveFileRef(t *testing.T) {
	tmp := t.TempDir()
	writeRefTestFile(t, filepath.Join(tmp, "internal", "stewref", "stewref.go"), "package stewref\n")
	ref, err := File("./internal/stewref/stewref.go")
	if err != nil {
		t.Fatalf("File() error = %v", err)
	}

	result, err := Resolve(ResolveOptions{
		TargetDir:     tmp,
		Ref:           ref,
		RequireExists: true,
	})
	if err != nil {
		t.Fatalf("Resolve() error = %v", err)
	}
	if result.Ref.String() != "file:internal/stewref/stewref.go" {
		t.Fatalf("Ref = %q", result.Ref.String())
	}
	if result.RelPath != "internal/stewref/stewref.go" {
		t.Fatalf("RelPath = %q", result.RelPath)
	}
	if result.AbsPath != filepath.Join(tmp, "internal", "stewref", "stewref.go") {
		t.Fatalf("AbsPath = %q", result.AbsPath)
	}
}

func TestResolveEntryRef(t *testing.T) {
	tmp := setupRefLedger(t, "decisions")
	entryName := "2026-06-20T191722Z-5cxsdb-use-ids-for-generated-entry-filenames.md"
	writeRefTestFile(t, filepath.Join(tmp, ".stew", "ledgers", "decisions", entryName), "entry\n")
	ref, err := Entry("decisions", entryName)
	if err != nil {
		t.Fatalf("Entry() error = %v", err)
	}

	result, err := Resolve(ResolveOptions{
		TargetDir:     tmp,
		Ref:           ref,
		RequireExists: true,
	})
	if err != nil {
		t.Fatalf("Resolve() error = %v", err)
	}
	if result.Ref.String() != "entry:decisions/"+entryName {
		t.Fatalf("Ref = %q", result.Ref.String())
	}
	if result.RelPath != ".stew/ledgers/decisions/"+entryName {
		t.Fatalf("RelPath = %q", result.RelPath)
	}
	if result.AbsPath != filepath.Join(tmp, ".stew", "ledgers", "decisions", entryName) {
		t.Fatalf("AbsPath = %q", result.AbsPath)
	}
}

func TestResolveRejectsUnknownLedger(t *testing.T) {
	tmp := t.TempDir()
	if err := os.MkdirAll(filepath.Join(tmp, ".stew"), 0o755); err != nil {
		t.Fatalf("mkdir .stew: %v", err)
	}
	ref, err := Entry("decisions", "2026-06-20T191722Z-5cxsdb-use-ids-for-generated-entry-filenames.md")
	if err != nil {
		t.Fatalf("Entry() error = %v", err)
	}

	_, err = Resolve(ResolveOptions{TargetDir: tmp, Ref: ref})
	if !errors.Is(err, ErrUnknownLedger) {
		t.Fatalf("Resolve() error = %v, want ErrUnknownLedger", err)
	}
}

func TestResolveRequireExistsRejectsMissingTarget(t *testing.T) {
	tmp := t.TempDir()
	ref, err := File("internal/missing.go")
	if err != nil {
		t.Fatalf("File() error = %v", err)
	}

	_, err = Resolve(ResolveOptions{
		TargetDir:     tmp,
		Ref:           ref,
		RequireExists: true,
	})
	if !errors.Is(err, ErrMissingTarget) {
		t.Fatalf("Resolve() error = %v, want ErrMissingTarget", err)
	}
}

func TestResolveAllowsMissingTargetWithoutRequireExists(t *testing.T) {
	tmp := t.TempDir()
	ref, err := File("internal/missing.go")
	if err != nil {
		t.Fatalf("File() error = %v", err)
	}

	result, err := Resolve(ResolveOptions{TargetDir: tmp, Ref: ref})
	if err != nil {
		t.Fatalf("Resolve() error = %v", err)
	}
	if result.RelPath != "internal/missing.go" {
		t.Fatalf("RelPath = %q", result.RelPath)
	}
}

func TestResolveRequireExistsRejectsDirectory(t *testing.T) {
	tmp := t.TempDir()
	if err := os.MkdirAll(filepath.Join(tmp, "internal"), 0o755); err != nil {
		t.Fatalf("mkdir internal: %v", err)
	}
	ref, err := File("internal")
	if err != nil {
		t.Fatalf("File() error = %v", err)
	}

	_, err = Resolve(ResolveOptions{
		TargetDir:     tmp,
		Ref:           ref,
		RequireExists: true,
	})
	if !errors.Is(err, ErrMissingTarget) {
		t.Fatalf("Resolve() error = %v, want ErrMissingTarget", err)
	}
}

func TestResolveRejectsLedgerEntryFileRef(t *testing.T) {
	tmp := t.TempDir()
	ref, err := File(".stew/ledgers/decisions/entry.md")
	if err != nil {
		t.Fatalf("File() error = %v", err)
	}

	_, err = Resolve(ResolveOptions{TargetDir: tmp, Ref: ref})
	if !errors.Is(err, ErrInvalidRef) {
		t.Fatalf("Resolve() error = %v, want ErrInvalidRef", err)
	}
}

func TestResolveRejectsManuallyInvalidRef(t *testing.T) {
	tmp := t.TempDir()

	_, err := Resolve(ResolveOptions{
		TargetDir: tmp,
		Ref: Ref{
			Kind:    KindFile,
			Payload: "../outside.go",
		},
	})
	if !errors.Is(err, ErrInvalidRef) {
		t.Fatalf("Resolve() error = %v, want ErrInvalidRef", err)
	}
}

func setupRefLedger(t *testing.T, ledger string) string {
	t.Helper()
	tmp := t.TempDir()
	if err := os.MkdirAll(filepath.Join(tmp, ".stew", "ledgers", ledger), 0o755); err != nil {
		t.Fatalf("mkdir ledger storage: %v", err)
	}
	if err := os.WriteFile(filepath.Join(tmp, ".stew", ledger+".spec.md"), []byte("# Spec\n"), 0o644); err != nil {
		t.Fatalf("write spec: %v", err)
	}
	return tmp
}

func writeRefTestFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("mkdir %s: %v", filepath.Dir(path), err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}
