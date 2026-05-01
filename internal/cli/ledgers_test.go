package cli

import (
	"bytes"
	"encoding/json"
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

	want := "alpha  Alpha description.\n" +
		"zeta   Zeta description.\n"
	if out.String() != want {
		t.Fatalf("stdout = %q, want %q", out.String(), want)
	}
}

func TestLedgersCommandOutputsJSON(t *testing.T) {
	tmp := setupCLIStewDir(t)
	chdirForTest(t, tmp)
	writeCLIStewFile(t, tmp, "stew.spec.md", "# Stew Spec\n\nShared model.\n")
	writeCLIStewFile(t, tmp, "zeta.spec.md", "# Zeta Spec\n\nZeta description.\n")
	writeCLIStewFile(t, tmp, "alpha.spec.md", "# Alpha Spec\n\nAlpha description.\n")
	var out bytes.Buffer

	err := ExecuteWithIO([]string{"ledgers", "--json"}, strings.NewReader(""), &out, &bytes.Buffer{})
	if err != nil {
		t.Fatalf("ExecuteWithIO() error = %v", err)
	}

	ledgers := decodeLedgersJSON(t, out.String())
	assertLedgerJSON(t, ledgers, []map[string]string{
		{"name": "alpha", "description": "Alpha description."},
		{"name": "zeta", "description": "Zeta description."},
	})
	assertNoPathFields(t, out.String(), tmp)
}

func TestLedgersCommandAcceptsPathFlag(t *testing.T) {
	tmp := setupCLIStewDir(t)
	writeCLIStewFile(t, tmp, "plans.spec.md", "# Plans Spec\n\nPlan records.\n")
	var out bytes.Buffer

	err := ExecuteWithIO([]string{"ledgers", "--path", tmp}, strings.NewReader(""), &out, &bytes.Buffer{})
	if err != nil {
		t.Fatalf("ExecuteWithIO() error = %v", err)
	}

	want := "plans  Plan records.\n"
	if out.String() != want {
		t.Fatalf("stdout = %q, want %q", out.String(), want)
	}
}

func TestLedgersCommandOutputsJSONWithPathFlag(t *testing.T) {
	tmp := setupCLIStewDir(t)
	writeCLIStewFile(t, tmp, "plans.spec.md", "# Plans Spec\n\nPlan records.\n")
	var out bytes.Buffer

	err := ExecuteWithIO([]string{"ledgers", "--json", "--path", tmp}, strings.NewReader(""), &out, &bytes.Buffer{})
	if err != nil {
		t.Fatalf("ExecuteWithIO() error = %v", err)
	}

	ledgers := decodeLedgersJSON(t, out.String())
	assertLedgerJSON(t, ledgers, []map[string]string{
		{"name": "plans", "description": "Plan records."},
	})
	assertNoPathFields(t, out.String(), tmp)
}

func TestLedgersHelpDocumentsDiscovery(t *testing.T) {
	output := executeHelp(t, "ledgers", "--help")

	required := []string{
		"List available Stew ledgers.",
		"discovers writable ledgers from ledger specs",
		"ledger name and description",
		"stew ledgers --json",
		"--json",
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

func decodeLedgersJSON(t *testing.T, output string) []map[string]string {
	t.Helper()
	if !strings.HasSuffix(output, "\n") {
		t.Fatalf("JSON output should end with newline, got %q", output)
	}
	if !json.Valid([]byte(output)) {
		t.Fatalf("stdout is not valid JSON: %q", output)
	}

	var topLevel map[string]json.RawMessage
	if err := json.Unmarshal([]byte(output), &topLevel); err != nil {
		t.Fatalf("unmarshal top-level JSON: %v", err)
	}
	if len(topLevel) != 1 {
		t.Fatalf("top-level keys = %v, want only ledgers", keys(topLevel))
	}
	rawLedgers, ok := topLevel["ledgers"]
	if !ok {
		t.Fatalf("top-level JSON missing ledgers key: %v", keys(topLevel))
	}

	var ledgers []map[string]string
	if err := json.Unmarshal(rawLedgers, &ledgers); err != nil {
		t.Fatalf("unmarshal ledgers JSON: %v", err)
	}
	return ledgers
}

func assertLedgerJSON(t *testing.T, got, want []map[string]string) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("ledgers len = %d, want %d", len(got), len(want))
	}
	for i := range want {
		if len(got[i]) != 2 {
			t.Fatalf("ledgers[%d] keys = %v, want only name and description", i, keys(got[i]))
		}
		if got[i]["name"] != want[i]["name"] || got[i]["description"] != want[i]["description"] {
			t.Fatalf("ledgers[%d] = %#v, want %#v", i, got[i], want[i])
		}
	}
}

func assertNoPathFields(t *testing.T, output, targetPath string) {
	t.Helper()
	for _, forbidden := range []string{".stew", "ledgerPath", "targetDir", targetPath} {
		if strings.Contains(output, forbidden) {
			t.Fatalf("JSON output should not expose %q\n--- output ---\n%s", forbidden, output)
		}
	}
}

func keys[V any](values map[string]V) []string {
	result := make([]string, 0, len(values))
	for key := range values {
		result = append(result, key)
	}
	return result
}

func writeCLIStewFile(t *testing.T, dir, name, body string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(dir, ".stew", name), []byte(body), 0o644); err != nil {
		t.Fatalf("write %s: %v", name, err)
	}
}
