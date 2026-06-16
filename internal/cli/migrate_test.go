package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestMigrateAtomicEntriesCommandMigratesLegacyLedger(t *testing.T) {
	tmp := setupCLILegacyLedger(t)
	var out bytes.Buffer

	err := ExecuteWithIO([]string{"migrate", "atomic-entries", "--path", tmp}, strings.NewReader(""), &out, &bytes.Buffer{})
	if err != nil {
		t.Fatalf("ExecuteWithIO() error = %v", err)
	}
	if out.String() != "Migrated 1 ledger(s), wrote 1 entries, and removed 1 legacy ledger file(s)\n" {
		t.Fatalf("stdout = %q", out.String())
	}
	if _, err := os.Stat(filepath.Join(tmp, ".stew", "iterations.md")); !os.IsNotExist(err) {
		t.Fatalf("legacy ledger should be removed, stat err = %v", err)
	}
	entry := readCLIFile(t, filepath.Join(tmp, ".stew", "ledgers", "iterations", "2026-05-01T000001Z-first-entry.md"))
	if !strings.Contains(entry, "Body") {
		t.Fatalf("entry content = %s", entry)
	}
}

func TestMigrateAtomicEntriesCommandDryRun(t *testing.T) {
	tmp := setupCLILegacyLedger(t)
	var out bytes.Buffer

	err := ExecuteWithIO([]string{"migrate", "atomic-entries", "--dry-run", "--path", tmp}, strings.NewReader(""), &out, &bytes.Buffer{})
	if err != nil {
		t.Fatalf("ExecuteWithIO() error = %v", err)
	}
	if out.String() != "Would migrate 1 ledger(s), write 1 entries, and remove 1 legacy ledger file(s)\n" {
		t.Fatalf("stdout = %q", out.String())
	}
	if _, err := os.Stat(filepath.Join(tmp, ".stew", "iterations.md")); err != nil {
		t.Fatalf("legacy ledger should remain, stat err = %v", err)
	}
}

func TestMigrateHelpDocumentsAtomicEntries(t *testing.T) {
	output := executeHelp(t, "migrate", "atomic-entries", "--help")
	required := []string{
		"Split monolithic Stew ledgers into atomic entry files.",
		"stew migrate atomic-entries --dry-run",
		"--dry-run",
		"--path string",
	}
	assertHelpContains(t, output, required)

	helpOutput := executeHelp(t, "help", "migrate", "atomic-entries")
	assertHelpContains(t, helpOutput, required)
}

func setupCLILegacyLedger(t *testing.T) string {
	t.Helper()
	tmp := t.TempDir()
	stewDir := filepath.Join(tmp, ".stew")
	if err := os.MkdirAll(stewDir, 0o755); err != nil {
		t.Fatalf("mkdir .stew: %v", err)
	}
	if err := os.WriteFile(filepath.Join(stewDir, "iterations.spec.md"), []byte("# Iterations Spec\n\nDescription.\n"), 0o644); err != nil {
		t.Fatalf("write spec: %v", err)
	}
	content := "# Iterations\n\n<!-- Managed by stew -->\n\n" +
		"## 2026-05-01T00:00:01Z — First entry\n\n" +
		"**Prompt:** Test prompt\n\n" +
		"Body\n\n---\n"
	if err := os.WriteFile(filepath.Join(stewDir, "iterations.md"), []byte(content), 0o644); err != nil {
		t.Fatalf("write legacy ledger: %v", err)
	}
	return tmp
}
