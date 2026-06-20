package stewmigrateatomic

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunMigratesLegacyLedgersAndRemovesOldFiles(t *testing.T) {
	tmp := setupLegacyRepo(t)
	writeLegacyLedger(t, tmp, "iterations", legacyLedgerContent(
		legacyEntry("2026-05-01T00:00:01Z", "First entry", "first body"),
		legacyEntry("2026-05-01T00:00:02Z", "Second entry", "second body"),
	))
	writeLegacyLedger(t, tmp, "decisions", legacyLedgerContent(
		legacyEntry("2026-05-01T00:00:03Z", "Use atomic entries", "decision body"),
	))

	result, err := Run(Options{TargetDir: tmp})
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if result.LedgerCount != 2 || result.EntryCount != 3 || result.RemovedCount != 2 {
		t.Fatalf("result = %#v, want 2 ledgers, 3 entries, 2 removed", result)
	}

	if _, err := os.Stat(filepath.Join(tmp, ".stew", "iterations.md")); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("iterations.md should be removed, stat err = %v", err)
	}
	first := readTestFile(t, filepath.Join(tmp, ".stew", "ledgers", "iterations", "2026-05-01T000001Z-first-entry.md"))
	if !strings.Contains(first, "first body") {
		t.Fatalf("first entry not written correctly: %s", first)
	}
	second := readTestFile(t, filepath.Join(tmp, ".stew", "ledgers", "iterations", "2026-05-01T000002Z-second-entry.md"))
	if !strings.Contains(second, "second body") {
		t.Fatalf("second entry not written correctly: %s", second)
	}
}

func TestRunDryRunDoesNotWriteOrRemove(t *testing.T) {
	tmp := setupLegacyRepo(t)
	writeLegacyLedger(t, tmp, "iterations", legacyLedgerContent(
		legacyEntry("2026-05-01T00:00:01Z", "First entry", "first body"),
	))

	result, err := Run(Options{TargetDir: tmp, DryRun: true})
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if result.WouldMigrate != 1 || result.WouldWrite != 1 || result.WouldRemove != 1 {
		t.Fatalf("dry run result = %#v", result)
	}
	if _, err := os.Stat(filepath.Join(tmp, ".stew", "iterations.md")); err != nil {
		t.Fatalf("legacy ledger should remain after dry run: %v", err)
	}
	if _, err := os.Stat(filepath.Join(tmp, ".stew", "ledgers", "iterations")); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("entry storage should not be created during dry run, stat err = %v", err)
	}
}

func TestRunMigratesEmptyLegacyLedger(t *testing.T) {
	tmp := setupLegacyRepo(t)
	writeLegacyLedger(t, tmp, "iterations", "# Iterations\n\n<!-- Managed by stew -->\n")

	result, err := Run(Options{TargetDir: tmp})
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if result.LedgerCount != 1 || result.EntryCount != 0 || result.RemovedCount != 1 {
		t.Fatalf("result = %#v, want empty migrated ledger", result)
	}
	info, err := os.Stat(filepath.Join(tmp, ".stew", "ledgers", "iterations"))
	if err != nil {
		t.Fatalf("stat entry storage: %v", err)
	}
	if !info.IsDir() {
		t.Fatalf("entry storage should be a directory")
	}
}

func TestRunRejectsMalformedLegacyLedgerWithoutEntries(t *testing.T) {
	tmp := setupLegacyRepo(t)
	writeLegacyLedger(t, tmp, "iterations", "# Iterations\n\nnot a valid entry\n")

	_, err := Run(Options{TargetDir: tmp})
	if err == nil {
		t.Fatalf("expected malformed legacy ledger error")
	}
	if !strings.Contains(err.Error(), "no valid Stew entries") {
		t.Fatalf("error = %v, want no valid entries", err)
	}
	if _, err := os.Stat(filepath.Join(tmp, ".stew", "iterations.md")); err != nil {
		t.Fatalf("legacy ledger should remain after error: %v", err)
	}
}

func TestRunDoesNotRemoveLegacyLedgerWhenAtomicWriteFails(t *testing.T) {
	tmp := setupLegacyRepo(t)
	writeLegacyLedger(t, tmp, "iterations", legacyLedgerContent(
		legacyEntry("2026-05-01T00:00:01Z", "First entry", "first body"),
	))
	entryDir := filepath.Join(tmp, ".stew", "ledgers", "iterations")
	if err := os.MkdirAll(entryDir, 0o755); err != nil {
		t.Fatalf("mkdir entry storage: %v", err)
	}
	if err := os.WriteFile(filepath.Join(entryDir, "existing.md"), []byte("existing"), 0o644); err != nil {
		t.Fatalf("write existing entry: %v", err)
	}

	_, err := Run(Options{TargetDir: tmp})
	if err == nil {
		t.Fatalf("expected write failure")
	}
	if !strings.Contains(err.Error(), "already contains markdown entries") {
		t.Fatalf("error = %v, want existing entries error", err)
	}
	if _, err := os.Stat(filepath.Join(tmp, ".stew", "iterations.md")); err != nil {
		t.Fatalf("legacy ledger should remain after write failure: %v", err)
	}
}

func TestRunSkipsLedgersWithoutLegacyFiles(t *testing.T) {
	tmp := setupLegacyRepo(t)

	result, err := Run(Options{TargetDir: tmp})
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if result.LedgerCount != 0 || result.SkippedCount != 1 {
		t.Fatalf("result = %#v, want skipped legacy-free ledger", result)
	}
}

func setupLegacyRepo(t *testing.T) string {
	t.Helper()
	tmp := t.TempDir()
	stewDir := filepath.Join(tmp, ".stew")
	if err := os.MkdirAll(stewDir, 0o755); err != nil {
		t.Fatalf("mkdir .stew: %v", err)
	}
	if err := os.WriteFile(filepath.Join(stewDir, "iterations.spec.md"), []byte("# Iterations Spec\n\nIterations.\n"), 0o644); err != nil {
		t.Fatalf("write spec: %v", err)
	}
	return tmp
}

func writeLegacyLedger(t *testing.T, dir, ledger, content string) {
	t.Helper()
	stewDir := filepath.Join(dir, ".stew")
	if err := os.WriteFile(filepath.Join(stewDir, ledger+".spec.md"), []byte("# "+ledger+" Spec\n\nDescription.\n"), 0o644); err != nil {
		t.Fatalf("write spec: %v", err)
	}
	if err := os.WriteFile(filepath.Join(stewDir, ledger+".md"), []byte(content), 0o644); err != nil {
		t.Fatalf("write legacy ledger: %v", err)
	}
}

func legacyLedgerContent(entries ...string) string {
	return "# Ledger\n\n<!-- Managed by stew -->\n\n" + strings.Join(entries, "\n")
}

func legacyEntry(timestamp, summary, body string) string {
	return "## " + timestamp + " — " + summary + "\n\n" +
		"**Prompt:** Test prompt\n\n" +
		body + "\n\n---\n"
}

func readTestFile(t *testing.T, path string) string {
	t.Helper()
	bytes, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return string(bytes)
}
