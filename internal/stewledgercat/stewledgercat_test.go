package stewledgercat

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestRunReadsExactLedgerContent(t *testing.T) {
	content := "## 2026-05-01T00:00:01Z — Entry\n\n**Prompt:** Test\n\nBody\n\n---\n"
	tmp := setupLedger(t, "iterations", content)

	result, err := Run(Options{
		TargetDir: tmp,
		Ledger:    "iterations",
	})
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if result.TargetDir != tmp {
		t.Fatalf("TargetDir = %q, want %q", result.TargetDir, tmp)
	}
	if result.LedgerPath != filepath.Join(".stew", "ledgers", "iterations") {
		t.Fatalf("LedgerPath = %q", result.LedgerPath)
	}
	want := content
	if result.Content != want {
		t.Fatalf("Content = %q, want %q", result.Content, want)
	}
}

func TestRunSupportsDefaultTargetDir(t *testing.T) {
	content := "## 2026-05-01T00:00:01Z — Decision\n\n**Prompt:** Test\n\nBody\n\n---\n"
	tmp := setupLedger(t, "decisions", content)
	chdirForLedgerCatTest(t, tmp)

	result, err := Run(Options{Ledger: "decisions"})
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if result.Content != content {
		t.Fatalf("Content = %q", result.Content)
	}
}

func TestRunRejectsUnknownLedger(t *testing.T) {
	tmp := setupStewDir(t)

	_, err := Run(Options{TargetDir: tmp, Ledger: "missing"})
	if !errors.Is(err, ErrUnknownLedger) {
		t.Fatalf("error = %v, want ErrUnknownLedger", err)
	}
}

func TestRunRejectsMissingLedgerFile(t *testing.T) {
	tmp := setupStewDir(t)
	writeFile(t, filepath.Join(tmp, ".stew", "plans.spec.md"), "# Plans Spec\n")

	_, err := Run(Options{TargetDir: tmp, Ledger: "plans"})
	if !errors.Is(err, ErrMissingLedger) {
		t.Fatalf("error = %v, want ErrMissingLedger", err)
	}
}

func TestRunRejectsInvalidLedgerNames(t *testing.T) {
	tmp := setupLedger(t, "iterations", "")
	cases := []string{
		"",
		" iterations",
		"iterations ",
		"stew",
		"../iterations",
		"iterations/notes",
		`iterations\notes`,
		"..",
		"iterations..old",
	}

	for _, name := range cases {
		t.Run(name, func(t *testing.T) {
			_, err := Run(Options{TargetDir: tmp, Ledger: name})
			if !errors.Is(err, ErrInvalidLedger) {
				t.Fatalf("error = %v, want ErrInvalidLedger", err)
			}
		})
	}
}

func TestRunRejectsDirectoryPaths(t *testing.T) {
	t.Run("spec directory", func(t *testing.T) {
		tmp := setupStewDir(t)
		if err := os.Mkdir(filepath.Join(tmp, ".stew", "plans.spec.md"), 0o755); err != nil {
			t.Fatalf("mkdir spec dir: %v", err)
		}

		_, err := Run(Options{TargetDir: tmp, Ledger: "plans"})
		if !errors.Is(err, ErrUnknownLedger) {
			t.Fatalf("error = %v, want ErrUnknownLedger", err)
		}
	})

	t.Run("ledger storage file", func(t *testing.T) {
		tmp := setupStewDir(t)
		writeFile(t, filepath.Join(tmp, ".stew", "plans.spec.md"), "# Plans Spec\n")
		if err := os.MkdirAll(filepath.Join(tmp, ".stew", "ledgers"), 0o755); err != nil {
			t.Fatalf("mkdir ledgers dir: %v", err)
		}
		if err := os.WriteFile(filepath.Join(tmp, ".stew", "ledgers", "plans"), []byte("not a dir"), 0o644); err != nil {
			t.Fatalf("write ledger storage file: %v", err)
		}

		_, err := Run(Options{TargetDir: tmp, Ledger: "plans"})
		if !errors.Is(err, ErrMissingLedger) {
			t.Fatalf("error = %v, want ErrMissingLedger", err)
		}
	})
}

func setupLedger(t *testing.T, ledger, content string) string {
	t.Helper()
	tmp := setupStewDir(t)
	writeFile(t, filepath.Join(tmp, ".stew", ledger+".spec.md"), "# Spec\n")
	entryDir := filepath.Join(tmp, ".stew", "ledgers", ledger)
	if err := os.MkdirAll(entryDir, 0o755); err != nil {
		t.Fatalf("mkdir ledger storage: %v", err)
	}
	if content != "" {
		writeFile(t, filepath.Join(entryDir, "2026-05-01T000001Z-entry.md"), content)
	}
	return tmp
}

func setupStewDir(t *testing.T) string {
	t.Helper()
	tmp := t.TempDir()
	if err := os.MkdirAll(filepath.Join(tmp, ".stew"), 0o755); err != nil {
		t.Fatalf("mkdir .stew: %v", err)
	}
	return tmp
}

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}

func chdirForLedgerCatTest(t *testing.T, dir string) {
	t.Helper()
	previous, err := os.Getwd()
	if err != nil {
		t.Fatalf("get working directory: %v", err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("chdir %s: %v", dir, err)
	}
	t.Cleanup(func() {
		if err := os.Chdir(previous); err != nil {
			t.Fatalf("restore working directory: %v", err)
		}
	})
}
