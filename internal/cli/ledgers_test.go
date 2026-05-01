package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLedgersCommandListsDiscoveredLedgers(t *testing.T) {
	tmp := setupCLIStewDir(t)
	chdirForTest(t, tmp)
	writeCLIStewFile(t, tmp, "stew.spec.md", "# Stew Spec\n\nShared model.\n")
	writeCLIStewFile(t, tmp, "zeta.spec.md", "# Zeta Spec\n\nZeta description.\n")
	writeCLIStewFile(t, tmp, "alpha.spec.md", "# Alpha Spec\n\nAlpha description.\n")
	var out bytes.Buffer

	err := ExecuteWithIO([]string{"ledgers"}, strings.NewReader(""), &out, &bytes.Buffer{})
	if err != nil {
		t.Fatalf("ExecuteWithIO() error = %v", err)
	}

	want := "alpha  " + filepath.Join(".stew", "alpha.md") + "  Alpha description.\n" +
		"zeta   " + filepath.Join(".stew", "zeta.md") + "   Zeta description.\n"
	if out.String() != want {
		t.Fatalf("stdout = %q, want %q", out.String(), want)
	}
}

func TestLedgersCommandAcceptsPathFlag(t *testing.T) {
	tmp := setupCLIStewDir(t)
	writeCLIStewFile(t, tmp, "plans.spec.md", "# Plans Spec\n\nPlan records.\n")
	var out bytes.Buffer

	err := ExecuteWithIO([]string{"ledgers", "--path", tmp}, strings.NewReader(""), &out, &bytes.Buffer{})
	if err != nil {
		t.Fatalf("ExecuteWithIO() error = %v", err)
	}

	want := "plans  " + filepath.Join(".stew", "plans.md") + "  Plan records.\n"
	if out.String() != want {
		t.Fatalf("stdout = %q, want %q", out.String(), want)
	}
}

func TestLedgersHelpDocumentsDiscovery(t *testing.T) {
	output := executeHelp(t, "ledgers", "--help")

	required := []string{
		"List available Stew ledgers.",
		"discovers writable ledgers from .stew/*.spec.md files",
		"excluding the",
		"stew ledgers --path /path/to/repo",
	}
	assertHelpContains(t, output, required)
}

func TestRootHelpIncludesLedgersCommand(t *testing.T) {
	output := executeHelp(t, "--help")

	required := []string{
		"ledgers",
		"List available ledgers",
	}
	assertHelpContains(t, output, required)
}

func writeCLIStewFile(t *testing.T, dir, name, body string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(dir, ".stew", name), []byte(body), 0o644); err != nil {
		t.Fatalf("write %s: %v", name, err)
	}
}
