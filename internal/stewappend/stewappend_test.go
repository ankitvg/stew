package stewappend

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/ankitvg/stew/internal/stewref"
)

func TestRunAppendsFormattedEntryFromMessage(t *testing.T) {
	tmp := setupLedger(t, "iterations")

	result, err := Run(Options{
		TargetDir:  tmp,
		Ledger:     "iterations",
		Prompt:     "Implement append",
		Summary:    "Add append command",
		Message:    "Created the append implementation.\n",
		MessageSet: true,
		Now: func() time.Time {
			return time.Date(2026, 4, 26, 12, 5, 7, 123, time.FixedZone("UTC+1", 3600))
		},
		NewEntryID: staticEntryID("k7p3qx"),
	})
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if result.LedgerPath != filepath.Join(".stew", "ledgers", "iterations") {
		t.Fatalf("LedgerPath = %q, want .stew/ledgers/iterations", result.LedgerPath)
	}
	wantEntryPath := filepath.Join(".stew", "ledgers", "iterations", "2026-04-26T110507Z-k7p3qx-add-append-command.md")
	if result.EntryPath != wantEntryPath {
		t.Fatalf("EntryPath = %q, want %q", result.EntryPath, wantEntryPath)
	}
	wantEntryRef := "entry:iterations/2026-04-26T110507Z-k7p3qx-add-append-command.md"
	if result.EntryRef != wantEntryRef {
		t.Fatalf("EntryRef = %q, want %q", result.EntryRef, wantEntryRef)
	}

	got := readFile(t, filepath.Join(tmp, result.EntryPath))
	want := "## 2026-04-26T11:05:07Z — Add append command\n\n" +
		"**Prompt:** Implement append\n\n" +
		"Created the append implementation.\n\n" +
		"---\n"
	if got != want {
		t.Fatalf("entry content mismatch\n--- got ---\n%s\n--- want ---\n%s", got, want)
	}
}

func TestRunAllowsMultilinePrompt(t *testing.T) {
	tmp := setupLedger(t, "decisions")

	result, err := Run(Options{
		TargetDir:  tmp,
		Ledger:     "decisions",
		Prompt:     "First line\nSecond line",
		Summary:    "Record prompt handling",
		Message:    "Decision body",
		MessageSet: true,
		Now:        fixedNow,
	})
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}

	got := readFile(t, filepath.Join(tmp, result.EntryPath))
	if !strings.Contains(got, "**Prompt:**\nFirst line\nSecond line\n\nDecision body") {
		t.Fatalf("multi-line prompt was not rendered verbatim: %s", got)
	}
}

func TestRunRejectsMissingLedgerStorageWhenSpecExists(t *testing.T) {
	tmp := setupLedger(t, "notes")
	if err := os.Remove(filepath.Join(tmp, ".stew", "ledgers", "notes")); err != nil {
		t.Fatalf("remove notes storage: %v", err)
	}

	_, err := Run(Options{
		TargetDir:  tmp,
		Ledger:     "notes",
		Prompt:     "Add note",
		Summary:    "Create notes entry",
		Message:    "Body",
		MessageSet: true,
		Now:        fixedNow,
	})
	if !errors.Is(err, ErrMissingStorage) {
		t.Fatalf("error = %v, want ErrMissingStorage", err)
	}
}

func TestRunReadsBodyFromFile(t *testing.T) {
	tmp := setupLedger(t, "iterations")
	bodyPath := filepath.Join(tmp, "body.md")
	if err := os.WriteFile(bodyPath, []byte("File body\n"), 0o644); err != nil {
		t.Fatalf("write body file: %v", err)
	}

	result, err := Run(Options{
		TargetDir: tmp,
		Ledger:    "iterations",
		Prompt:    "Use file body",
		Summary:   "Read file body",
		FilePath:  bodyPath,
		Now:       fixedNow,
	})
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}

	got := readFile(t, filepath.Join(tmp, result.EntryPath))
	if !strings.Contains(got, "\n\nFile body\n\n---\n") {
		t.Fatalf("file body missing from entry: %s", got)
	}
}

func TestRunReadsBodyFromStdin(t *testing.T) {
	tmp := setupLedger(t, "iterations")

	result, err := Run(Options{
		TargetDir: tmp,
		Ledger:    "iterations",
		Prompt:    "Use stdin body",
		Summary:   "Read stdin body",
		Stdin:     strings.NewReader("stdin body\n"),
		StdinIsTTY: func() bool {
			return false
		},
		Now: fixedNow,
	})
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}

	got := readFile(t, filepath.Join(tmp, result.EntryPath))
	if !strings.Contains(got, "\n\nstdin body\n\n---\n") {
		t.Fatalf("stdin body missing from entry: %s", got)
	}
}

func TestRunUsesNumericSuffixForSameSecondCollision(t *testing.T) {
	tmp := setupLedger(t, "iterations")
	opts := Options{
		TargetDir:  tmp,
		Ledger:     "iterations",
		Prompt:     "Prompt",
		Summary:    "Same summary",
		Message:    "Body",
		MessageSet: true,
		Now:        fixedNow,
		NewEntryID: staticEntryID("abc234"),
	}
	first, err := Run(opts)
	if err != nil {
		t.Fatalf("first Run() error = %v", err)
	}
	second, err := Run(opts)
	if err != nil {
		t.Fatalf("second Run() error = %v", err)
	}

	if filepath.Base(first.EntryPath) != "2026-04-26T183923Z-abc234-same-summary.md" {
		t.Fatalf("first EntryPath = %q", first.EntryPath)
	}
	if filepath.Base(second.EntryPath) != "2026-04-26T183923Z-abc234-same-summary-2.md" {
		t.Fatalf("second EntryPath = %q", second.EntryPath)
	}
}

func TestRunCreatesLinksForLinkedFiles(t *testing.T) {
	tmp := setupLedger(t, "iterations")
	writeAppendTestFile(t, filepath.Join(tmp, "internal", "cli", "ledger.go"), "package cli\n")
	writeAppendTestFile(t, filepath.Join(tmp, "internal", "stewappend", "stewappend.go"), "package stewappend\n")

	result, err := Run(Options{
		TargetDir:  tmp,
		Ledger:     "iterations",
		Prompt:     "Prompt",
		Summary:    "Link files",
		Message:    "Body",
		MessageSet: true,
		LinkFiles: []string{
			`./internal\cli/ledger.go`,
			"internal/stewappend/stewappend.go",
		},
		Now:        fixedNow,
		NewEntryID: staticEntryID("abc234"),
	})
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}

	if len(result.Links) != 2 {
		t.Fatalf("Links len = %d, want 2", len(result.Links))
	}
	wantEntryRef := "entry:iterations/2026-04-26T183923Z-abc234-link-files.md"
	if result.Links[0].Source != wantEntryRef || result.Links[0].Target != "file:internal/cli/ledger.go" || result.Links[0].CreatedAt != "2026-04-26T18:39:23Z" {
		t.Fatalf("Links[0] = %#v", result.Links[0])
	}
	if result.Links[1].Source != wantEntryRef || result.Links[1].Target != "file:internal/stewappend/stewappend.go" || result.Links[1].CreatedAt != "2026-04-26T18:39:23Z" {
		t.Fatalf("Links[1] = %#v", result.Links[1])
	}

	linkFiles, err := os.ReadDir(filepath.Join(tmp, ".stew", "links"))
	if err != nil {
		t.Fatalf("read links dir: %v", err)
	}
	if len(linkFiles) != 2 {
		t.Fatalf("link files len = %d, want 2", len(linkFiles))
	}
}

func TestRunDedupesLinkedFiles(t *testing.T) {
	tmp := setupLedger(t, "iterations")
	writeAppendTestFile(t, filepath.Join(tmp, "internal", "cli", "ledger.go"), "package cli\n")

	result, err := Run(Options{
		TargetDir:  tmp,
		Ledger:     "iterations",
		Prompt:     "Prompt",
		Summary:    "Link files",
		Message:    "Body",
		MessageSet: true,
		LinkFiles: []string{
			"internal/cli/ledger.go",
			"./internal/cli/ledger.go",
		},
		Now:        fixedNow,
		NewEntryID: staticEntryID("abc234"),
	})
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if len(result.Links) != 1 {
		t.Fatalf("Links len = %d, want 1", len(result.Links))
	}
}

func TestRunRejectsMissingLinkedFileBeforeCreatingEntry(t *testing.T) {
	tmp := setupLedger(t, "iterations")

	_, err := Run(Options{
		TargetDir:  tmp,
		Ledger:     "iterations",
		Prompt:     "Prompt",
		Summary:    "Missing link",
		Message:    "Body",
		MessageSet: true,
		LinkFiles:  []string{"internal/cli/missing.go"},
		Now:        fixedNow,
		NewEntryID: staticEntryID("abc234"),
	})
	if !errors.Is(err, stewref.ErrMissingTarget) {
		t.Fatalf("Run() error = %v, want ErrMissingTarget", err)
	}
	if got := entryFileCount(t, tmp, "iterations"); got != 0 {
		t.Fatalf("entry files len = %d, want 0", got)
	}
}

func TestRunKeepsEntryWhenLinkWriteFails(t *testing.T) {
	tmp := setupLedger(t, "iterations")
	writeAppendTestFile(t, filepath.Join(tmp, "internal", "cli", "ledger.go"), "package cli\n")
	if err := os.WriteFile(filepath.Join(tmp, ".stew", "links"), []byte("not a dir"), 0o644); err != nil {
		t.Fatalf("write links file: %v", err)
	}

	_, err := Run(Options{
		TargetDir:  tmp,
		Ledger:     "iterations",
		Prompt:     "Prompt",
		Summary:    "Link failure",
		Message:    "Body",
		MessageSet: true,
		LinkFiles:  []string{"internal/cli/ledger.go"},
		Now:        fixedNow,
		NewEntryID: staticEntryID("abc234"),
	})
	if err == nil {
		t.Fatalf("expected link write error")
	}
	if !strings.Contains(err.Error(), "create links for entry:iterations/2026-04-26T183923Z-abc234-link-failure.md") {
		t.Fatalf("error = %v, want created entry ref context", err)
	}
	if got := entryFileCount(t, tmp, "iterations"); got != 1 {
		t.Fatalf("entry files len = %d, want 1", got)
	}
}

func TestRunRejectsInvalidGeneratedEntryID(t *testing.T) {
	tmp := setupLedger(t, "iterations")

	_, err := Run(Options{
		TargetDir:  tmp,
		Ledger:     "iterations",
		Prompt:     "Prompt",
		Summary:    "Summary",
		Message:    "Body",
		MessageSet: true,
		Now:        fixedNow,
		NewEntryID: staticEntryID("bad-id"),
	})
	if err == nil {
		t.Fatalf("expected invalid entry id error")
	}
	if !strings.Contains(err.Error(), "entry id") {
		t.Fatalf("error = %v, want entry id context", err)
	}
}

func TestRunRejectsMissingAndConflictingBodySources(t *testing.T) {
	tmp := setupLedger(t, "iterations")

	_, err := Run(Options{
		TargetDir: tmp,
		Ledger:    "iterations",
		Prompt:    "Missing body",
		Summary:   "Reject missing body",
		Stdin:     strings.NewReader("ignored"),
		StdinIsTTY: func() bool {
			return true
		},
		Now: fixedNow,
	})
	if !errors.Is(err, ErrMissingBody) {
		t.Fatalf("error = %v, want ErrMissingBody", err)
	}

	_, err = Run(Options{
		TargetDir:  tmp,
		Ledger:     "iterations",
		Prompt:     "Conflicting body",
		Summary:    "Reject body conflict",
		Message:    "message",
		MessageSet: true,
		Stdin:      strings.NewReader("stdin"),
		StdinIsTTY: func() bool {
			return false
		},
		Now: fixedNow,
	})
	if !errors.Is(err, ErrBodyConflict) {
		t.Fatalf("error = %v, want ErrBodyConflict", err)
	}
}

func TestRunRejectsInvalidLedgers(t *testing.T) {
	tmp := setupLedger(t, "iterations")

	cases := []string{"", "stew", "../iterations", "nested/name", `nested\name`, "bad..name", " iterations"}
	for _, ledger := range cases {
		t.Run(ledger, func(t *testing.T) {
			_, err := Run(Options{
				TargetDir:  tmp,
				Ledger:     ledger,
				Prompt:     "Prompt",
				Summary:    "Summary",
				Message:    "Body",
				MessageSet: true,
				Now:        fixedNow,
			})
			if !errors.Is(err, ErrInvalidLedger) {
				t.Fatalf("error = %v, want ErrInvalidLedger", err)
			}
		})
	}
}

func TestRunRejectsUnknownLedger(t *testing.T) {
	tmp := setupLedger(t, "iterations")

	_, err := Run(Options{
		TargetDir:  tmp,
		Ledger:     "decisions",
		Prompt:     "Prompt",
		Summary:    "Summary",
		Message:    "Body",
		MessageSet: true,
		Now:        fixedNow,
	})
	if !errors.Is(err, ErrUnknownLedger) {
		t.Fatalf("error = %v, want ErrUnknownLedger", err)
	}
}

func TestRunRejectsMissingPromptAndSummary(t *testing.T) {
	tmp := setupLedger(t, "iterations")

	_, err := Run(Options{
		TargetDir:  tmp,
		Ledger:     "iterations",
		Summary:    "Summary",
		Message:    "Body",
		MessageSet: true,
		Now:        fixedNow,
	})
	if !errors.Is(err, ErrMissingPrompt) {
		t.Fatalf("error = %v, want ErrMissingPrompt", err)
	}

	_, err = Run(Options{
		TargetDir:  tmp,
		Ledger:     "iterations",
		Prompt:     "Prompt",
		Summary:    "bad\nsummary",
		Message:    "Body",
		MessageSet: true,
		Now:        fixedNow,
	})
	if !errors.Is(err, ErrMissingSummary) {
		t.Fatalf("error = %v, want ErrMissingSummary", err)
	}
}

func setupLedger(t *testing.T, ledger string) string {
	t.Helper()
	tmp := t.TempDir()
	stewDir := filepath.Join(tmp, ".stew")
	if err := os.MkdirAll(stewDir, 0o755); err != nil {
		t.Fatalf("mkdir .stew: %v", err)
	}
	if err := os.WriteFile(filepath.Join(stewDir, ledger+".spec.md"), []byte("# Spec\n"), 0o644); err != nil {
		t.Fatalf("write spec: %v", err)
	}
	if err := os.MkdirAll(filepath.Join(stewDir, "ledgers", ledger), 0o755); err != nil {
		t.Fatalf("mkdir ledger storage: %v", err)
	}
	return tmp
}

func readFile(t *testing.T, path string) string {
	t.Helper()
	bytes, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return string(bytes)
}

func fixedNow() time.Time {
	return time.Date(2026, 4, 26, 18, 39, 23, 987, time.UTC)
}

func staticEntryID(id string) func() (string, error) {
	return func() (string, error) {
		return id, nil
	}
}

func writeAppendTestFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("mkdir %s: %v", filepath.Dir(path), err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}

func entryFileCount(t *testing.T, dir, ledger string) int {
	t.Helper()
	entries, err := os.ReadDir(filepath.Join(dir, ".stew", "ledgers", ledger))
	if err != nil {
		t.Fatalf("read entries: %v", err)
	}
	return len(entries)
}
