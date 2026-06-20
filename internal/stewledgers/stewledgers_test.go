package stewledgers

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestListDiscoversWritableLedgers(t *testing.T) {
	tmp := setupStewDir(t)
	writeSpec(t, tmp, "stew.spec.md", "# Stew Spec\n\nShared model.\n")
	writeSpec(t, tmp, "zeta.spec.md", "# Zeta Spec\n\nZeta description.\n")
	writeSpec(t, tmp, "alpha.spec.md", "# Alpha Spec\n\nAlpha starts here\nand wraps across lines.\n\n## Body\n\nOther text.\n")
	writeSpec(t, tmp, "notes.md", "not a spec\n")

	result, err := List(Options{TargetDir: tmp})
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}
	if result.TargetDir != tmp {
		t.Fatalf("TargetDir = %q, want %q", result.TargetDir, tmp)
	}

	want := []Ledger{
		{
			Name:        "alpha",
			LedgerPath:  filepath.Join(".stew", "ledgers", "alpha"),
			Description: "Alpha starts here and wraps across lines.",
		},
		{
			Name:        "zeta",
			LedgerPath:  filepath.Join(".stew", "ledgers", "zeta"),
			Description: "Zeta description.",
		},
	}
	if len(result.Ledgers) != len(want) {
		t.Fatalf("Ledgers len = %d, want %d", len(result.Ledgers), len(want))
	}
	for i := range want {
		if result.Ledgers[i] != want[i] {
			t.Fatalf("Ledgers[%d] = %#v, want %#v", i, result.Ledgers[i], want[i])
		}
	}
}

func TestListErrorsWhenStewDirMissing(t *testing.T) {
	tmp := t.TempDir()

	_, err := List(Options{TargetDir: tmp})
	if !errors.Is(err, ErrMissingStewDir) {
		t.Fatalf("error = %v, want ErrMissingStewDir", err)
	}
}

func TestListErrorsWhenNoWritableLedgerSpecs(t *testing.T) {
	tmp := setupStewDir(t)
	writeSpec(t, tmp, "stew.spec.md", "# Stew Spec\n\nShared model.\n")
	writeSpec(t, tmp, "notes.md", "not a spec\n")

	_, err := List(Options{TargetDir: tmp})
	if !errors.Is(err, ErrNoLedgers) {
		t.Fatalf("error = %v, want ErrNoLedgers", err)
	}
}

func setupStewDir(t *testing.T) string {
	t.Helper()
	tmp := t.TempDir()
	if err := os.Mkdir(filepath.Join(tmp, ".stew"), 0o755); err != nil {
		t.Fatalf("mkdir .stew: %v", err)
	}
	return tmp
}

func writeSpec(t *testing.T, dir, name, body string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(dir, ".stew", name), []byte(body), 0o644); err != nil {
		t.Fatalf("write %s: %v", name, err)
	}
}
