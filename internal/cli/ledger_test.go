package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLedgerNewCommandCreatesLedgerAndSpec(t *testing.T) {
	tmp := setupCLIStewDir(t)
	chdirForTest(t, tmp)
	var out bytes.Buffer

	cmd := newRootCmd()
	cmd.SetOut(&out)
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetArgs([]string{
		"ledger", "new", "plans",
		"--description", "Reasoning artifacts for future work.",
		"--threshold", "Append when a plan captures durable intent or tradeoffs.",
	})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	wantOut := "Created " + filepath.Join(".stew", "plans.md") + "\n" +
		"Created " + filepath.Join(".stew", "plans.spec.md") + "\n"
	if out.String() != wantOut {
		t.Fatalf("stdout = %q, want %q", out.String(), wantOut)
	}

	ledger := readCLIFile(t, filepath.Join(tmp, ".stew", "plans.md"))
	if ledger != "# Plans\n\n<!-- Managed by stew -->\n" {
		t.Fatalf("ledger content = %q", ledger)
	}
	spec := readCLIFile(t, filepath.Join(tmp, ".stew", "plans.spec.md"))
	if !strings.Contains(spec, "Reasoning artifacts for future work.") {
		t.Fatalf("spec missing description: %s", spec)
	}
	if !strings.Contains(spec, "Append when a plan captures durable intent or tradeoffs.") {
		t.Fatalf("spec missing threshold: %s", spec)
	}

	appendCmd := newRootCmd()
	appendCmd.SetOut(&bytes.Buffer{})
	appendCmd.SetErr(&bytes.Buffer{})
	appendCmd.SetArgs([]string{
		"append", "plans",
		"--prompt", "Capture plan",
		"--summary", "Record plan",
		"-m", "Plan body",
	})
	if err := appendCmd.Execute(); err != nil {
		t.Fatalf("append Execute() error = %v", err)
	}

	fullSpecCmd := newRootCmd()
	var fullSpecOut bytes.Buffer
	fullSpecCmd.SetOut(&fullSpecOut)
	fullSpecCmd.SetErr(&bytes.Buffer{})
	fullSpecCmd.SetArgs([]string{"full-spec"})
	if err := fullSpecCmd.Execute(); err != nil {
		t.Fatalf("full-spec Execute() error = %v", err)
	}
	if !strings.Contains(fullSpecOut.String(), "<!-- .stew/plans.spec.md -->") {
		t.Fatalf("full-spec missing plans spec: %s", fullSpecOut.String())
	}
}

func TestLedgerNewCommandQuietSuppressesOutput(t *testing.T) {
	tmp := setupCLIStewDir(t)
	chdirForTest(t, tmp)
	var out bytes.Buffer

	cmd := newRootCmd()
	cmd.SetOut(&out)
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetArgs([]string{
		"ledger", "new", "plans",
		"--quiet",
	})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if out.String() != "" {
		t.Fatalf("stdout = %q, want empty", out.String())
	}
}

func setupCLIStewDir(t *testing.T) string {
	t.Helper()
	tmp := t.TempDir()
	if err := os.MkdirAll(filepath.Join(tmp, ".stew"), 0o755); err != nil {
		t.Fatalf("mkdir .stew: %v", err)
	}
	return tmp
}

func chdirForTest(t *testing.T, dir string) {
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
