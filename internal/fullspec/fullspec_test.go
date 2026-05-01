package fullspec

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadSortsStewSpecFirstAndConcats(t *testing.T) {
	tmp := t.TempDir()
	stewDir := filepath.Join(tmp, ".stew")
	if err := os.MkdirAll(stewDir, 0o755); err != nil {
		t.Fatalf("mkdir .stew: %v", err)
	}

	files := map[string]string{
		"zeta.spec.md":  "# Zeta\n",
		"stew.spec.md":  "# Stew\n",
		"alpha.spec.md": "# Alpha\n",
		"notes.md":      "not a spec\n",
	}
	for name, body := range files {
		if err := os.WriteFile(filepath.Join(stewDir, name), []byte(body), 0o644); err != nil {
			t.Fatalf("write %s: %v", name, err)
		}
	}

	result, err := Load(tmp)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	wantFiles := []string{
		filepath.Join(".stew", "stew.spec.md"),
		filepath.Join(".stew", "alpha.spec.md"),
		filepath.Join(".stew", "zeta.spec.md"),
	}
	if len(result.SpecFiles) != len(wantFiles) {
		t.Fatalf("SpecFiles len = %d, want %d", len(result.SpecFiles), len(wantFiles))
	}
	for i, want := range wantFiles {
		if got := result.SpecFiles[i]; got != want {
			t.Fatalf("SpecFiles[%d] = %q, want %q", i, got, want)
		}
	}

	stewIdx := strings.Index(result.Content, "<!-- stew spec -->")
	alphaIdx := strings.Index(result.Content, "<!-- alpha spec -->")
	zetaIdx := strings.Index(result.Content, "<!-- zeta spec -->")
	if stewIdx < 0 || alphaIdx < 0 || zetaIdx < 0 {
		t.Fatalf("missing expected source markers in output: %s", result.Content)
	}
	if !(stewIdx < alphaIdx && alphaIdx < zetaIdx) {
		t.Fatalf("unexpected source order in output: %s", result.Content)
	}
}

func TestLoadErrorsWhenStewDirMissing(t *testing.T) {
	tmp := t.TempDir()
	_, err := Load(tmp)
	if err == nil {
		t.Fatalf("expected error")
	}
	if !errors.Is(err, ErrMissingStewDir) {
		t.Fatalf("error = %v, want ErrMissingStewDir", err)
	}
}

func TestLoadErrorsWhenNoSpecFiles(t *testing.T) {
	tmp := t.TempDir()
	stewDir := filepath.Join(tmp, ".stew")
	if err := os.MkdirAll(stewDir, 0o755); err != nil {
		t.Fatalf("mkdir .stew: %v", err)
	}
	if err := os.WriteFile(filepath.Join(stewDir, "notes.md"), []byte("x\n"), 0o644); err != nil {
		t.Fatalf("write notes.md: %v", err)
	}

	_, err := Load(tmp)
	if err == nil {
		t.Fatalf("expected error")
	}
	if !errors.Is(err, ErrNoSpecFiles) {
		t.Fatalf("error = %v, want ErrNoSpecFiles", err)
	}
}
