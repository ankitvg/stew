package stewledger

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunCreatesLedgerAndSpec(t *testing.T) {
	tmp := setupStewDir(t)

	result, err := Run(Options{
		TargetDir:   tmp,
		Name:        "security-audits",
		Description: "Tracks security review findings and follow-up rationale.",
		Threshold:   "Append when a security review produces durable findings or tradeoffs.",
	})
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if result.TargetDir != tmp {
		t.Fatalf("TargetDir = %q, want %q", result.TargetDir, tmp)
	}
	if result.LedgerPath != filepath.Join(".stew", "ledgers", "security-audits") {
		t.Fatalf("LedgerPath = %q", result.LedgerPath)
	}
	if result.SpecPath != filepath.Join(".stew", "security-audits.spec.md") {
		t.Fatalf("SpecPath = %q", result.SpecPath)
	}

	ledgerInfo, err := os.Stat(filepath.Join(tmp, ".stew", "ledgers", "security-audits"))
	if err != nil {
		t.Fatalf("stat ledger storage: %v", err)
	}
	if !ledgerInfo.IsDir() {
		t.Fatalf("ledger storage should be a directory")
	}

	spec := readTestFile(t, filepath.Join(tmp, ".stew", "security-audits.spec.md"))
	required := []string{
		"# Security Audits Spec",
		"Tracks security review findings and follow-up rationale.",
		"## Body",
		"Entries use freeform markdown under the standard Stew entry header.",
		"## When to append",
		"Append when a security review produces durable findings or tradeoffs.",
	}
	for _, want := range required {
		if !strings.Contains(spec, want) {
			t.Fatalf("spec missing %q\n--- spec ---\n%s", want, spec)
		}
	}
}

func TestRunUsesTODOFallbacks(t *testing.T) {
	tmp := setupStewDir(t)

	_, err := Run(Options{
		TargetDir: tmp,
		Name:      "plans",
	})
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}

	spec := readTestFile(t, filepath.Join(tmp, ".stew", "plans.spec.md"))
	required := []string{
		"TODO: Describe what this ledger records and who should read it.",
		"TODO: Explain when an entry belongs in this ledger instead of another ledger.",
	}
	for _, want := range required {
		if !strings.Contains(spec, want) {
			t.Fatalf("spec missing TODO fallback %q\n--- spec ---\n%s", want, spec)
		}
	}
}

func TestRunRejectsInvalidAndReservedNames(t *testing.T) {
	tmp := setupStewDir(t)
	cases := []string{
		"",
		".plans",
		"Plans",
		"plan_notes",
		"plan notes",
		"plan--notes",
		"-plans",
		"plans-",
		"plans/review",
		"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
		"stew",
		"config",
		"stews",
		"ledger",
		"ledgers",
		"spec",
		"specs",
	}

	for _, name := range cases {
		t.Run(name, func(t *testing.T) {
			_, err := Run(Options{TargetDir: tmp, Name: name})
			if !errors.Is(err, ErrInvalidName) {
				t.Fatalf("error = %v, want ErrInvalidName", err)
			}
		})
	}
}

func TestRunRejectsMissingStewDirectory(t *testing.T) {
	tmp := t.TempDir()

	_, err := Run(Options{
		TargetDir: tmp,
		Name:      "plans",
	})
	if !errors.Is(err, ErrMissingStewDir) {
		t.Fatalf("error = %v, want ErrMissingStewDir", err)
	}
	if !strings.Contains(err.Error(), "run `stew init` first") {
		t.Fatalf("error missing init guidance: %v", err)
	}
}

func TestRunRejectsExistingFilesWithoutClobbering(t *testing.T) {
	t.Run("legacy ledger exists", func(t *testing.T) {
		tmp := setupStewDir(t)
		ledgerPath := filepath.Join(tmp, ".stew", "plans.md")
		original := "# Existing Plans\n"
		if err := os.WriteFile(ledgerPath, []byte(original), 0o644); err != nil {
			t.Fatalf("write existing ledger: %v", err)
		}

		_, err := Run(Options{TargetDir: tmp, Name: "plans"})
		if !errors.Is(err, ErrLedgerExists) {
			t.Fatalf("error = %v, want ErrLedgerExists", err)
		}
		if got := readTestFile(t, ledgerPath); got != original {
			t.Fatalf("existing ledger was clobbered: %q", got)
		}
		if _, err := os.Stat(filepath.Join(tmp, ".stew", "plans.spec.md")); !errors.Is(err, os.ErrNotExist) {
			t.Fatalf("spec should not be created after preflight failure, stat err = %v", err)
		}
	})

	t.Run("ledger storage exists", func(t *testing.T) {
		tmp := setupStewDir(t)
		ledgerPath := filepath.Join(tmp, ".stew", "ledgers", "plans")
		if err := os.MkdirAll(ledgerPath, 0o755); err != nil {
			t.Fatalf("mkdir existing ledger storage: %v", err)
		}

		_, err := Run(Options{TargetDir: tmp, Name: "plans"})
		if !errors.Is(err, ErrLedgerExists) {
			t.Fatalf("error = %v, want ErrLedgerExists", err)
		}
		if _, err := os.Stat(filepath.Join(tmp, ".stew", "plans.spec.md")); !errors.Is(err, os.ErrNotExist) {
			t.Fatalf("spec should not be created after preflight failure, stat err = %v", err)
		}
	})

	t.Run("spec exists", func(t *testing.T) {
		tmp := setupStewDir(t)
		specPath := filepath.Join(tmp, ".stew", "plans.spec.md")
		original := "# Existing Plans Spec\n"
		if err := os.WriteFile(specPath, []byte(original), 0o644); err != nil {
			t.Fatalf("write existing spec: %v", err)
		}

		_, err := Run(Options{TargetDir: tmp, Name: "plans"})
		if !errors.Is(err, ErrLedgerExists) {
			t.Fatalf("error = %v, want ErrLedgerExists", err)
		}
		if got := readTestFile(t, specPath); got != original {
			t.Fatalf("existing spec was clobbered: %q", got)
		}
		if _, err := os.Stat(filepath.Join(tmp, ".stew", "ledgers", "plans")); !errors.Is(err, os.ErrNotExist) {
			t.Fatalf("ledger storage should not be created after preflight failure, stat err = %v", err)
		}
	})
}

func setupStewDir(t *testing.T) string {
	t.Helper()
	tmp := t.TempDir()
	if err := os.Mkdir(filepath.Join(tmp, ".stew"), 0o755); err != nil {
		t.Fatalf("mkdir .stew: %v", err)
	}
	return tmp
}

func readTestFile(t *testing.T, path string) string {
	t.Helper()
	bytes, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return string(bytes)
}
