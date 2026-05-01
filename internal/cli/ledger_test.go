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

	err := ExecuteWithIO([]string{
		"ledger", "new", "plans",
		"--description", "Reasoning artifacts for future work.",
		"--threshold", "Append when a plan captures durable intent or tradeoffs.",
	}, strings.NewReader(""), &out, &bytes.Buffer{})
	if err != nil {
		t.Fatalf("ExecuteWithIO() error = %v", err)
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

	err = ExecuteWithIO([]string{
		"append", "plans",
		"--prompt", "Capture plan",
		"--summary", "Record plan",
		"-m", "Plan body",
	}, strings.NewReader(""), &bytes.Buffer{}, &bytes.Buffer{})
	if err != nil {
		t.Fatalf("append ExecuteWithIO() error = %v", err)
	}

	var fullSpecOut bytes.Buffer
	if err := ExecuteWithIO([]string{"full-spec"}, strings.NewReader(""), &fullSpecOut, &bytes.Buffer{}); err != nil {
		t.Fatalf("full-spec ExecuteWithIO() error = %v", err)
	}
	if !strings.Contains(fullSpecOut.String(), "<!-- .stew/plans.spec.md -->") {
		t.Fatalf("full-spec missing plans spec: %s", fullSpecOut.String())
	}
}

func TestLedgerNewCommandQuietSuppressesOutput(t *testing.T) {
	tmp := setupCLIStewDir(t)
	chdirForTest(t, tmp)
	var out bytes.Buffer

	err := ExecuteWithIO([]string{
		"ledger", "new", "plans",
		"--quiet",
	}, strings.NewReader(""), &out, &bytes.Buffer{})
	if err != nil {
		t.Fatalf("ExecuteWithIO() error = %v", err)
	}
	if out.String() != "" {
		t.Fatalf("stdout = %q, want empty", out.String())
	}
}

func TestLedgerCatCommandPrintsRawLedger(t *testing.T) {
	tmp := setupCLILedger(t, "iterations")
	chdirForTest(t, tmp)
	var out bytes.Buffer

	err := ExecuteWithIO([]string{"ledger", "cat", "iterations"}, strings.NewReader(""), &out, &bytes.Buffer{})
	if err != nil {
		t.Fatalf("ExecuteWithIO() error = %v", err)
	}

	want := "# Iterations\n\n<!-- Managed by stew -->\n"
	if out.String() != want {
		t.Fatalf("stdout = %q, want %q", out.String(), want)
	}
}

func TestLedgerCatCommandAcceptsPathFlag(t *testing.T) {
	tmp := setupCLILedger(t, "decisions")
	var out bytes.Buffer

	err := ExecuteWithIO([]string{"ledger", "cat", "decisions", "--path", tmp}, strings.NewReader(""), &out, &bytes.Buffer{})
	if err != nil {
		t.Fatalf("ExecuteWithIO() error = %v", err)
	}

	want := "# Iterations\n\n<!-- Managed by stew -->\n"
	if out.String() != want {
		t.Fatalf("stdout = %q, want %q", out.String(), want)
	}
}

func TestLedgerCatHelpDocumentsRawOutput(t *testing.T) {
	output := executeHelp(t, "ledger", "cat", "--help")

	required := []string{
		"Print a Stew ledger.",
		"raw .stew/<ledger>.md file to stdout",
		"stew ledger cat iterations | grep",
		"--path string",
	}
	assertHelpContains(t, output, required)

	helpOutput := executeHelp(t, "help", "ledger", "cat")
	assertHelpContains(t, helpOutput, required)
}

func TestLedgerHelpIncludesCatCommand(t *testing.T) {
	output := executeHelp(t, "ledger", "--help")

	required := []string{
		"cat",
		"Print a stew ledger",
		"new",
		"Create a custom stew ledger",
	}
	assertHelpContains(t, output, required)
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
