package stewledgertail

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunReturnsLastEntriesInFileOrder(t *testing.T) {
	content := ledgerContent(
		entry("2026-05-01T00:00:01Z", "First", "first body"),
		entry("2026-05-01T00:00:02Z", "Second", "second body"),
		entry("2026-05-01T00:00:03Z", "Third", "third body"),
	)
	tmp := setupTailLedger(t, "iterations", content)

	result, err := Run(Options{
		TargetDir: tmp,
		Ledger:    "iterations",
		Limit:     2,
	})
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}

	want := entry("2026-05-01T00:00:02Z", "Second", "second body") + "\n" +
		entry("2026-05-01T00:00:03Z", "Third", "third body")
	if result.Content != want {
		t.Fatalf("Content = %q, want %q", result.Content, want)
	}
	if strings.Contains(result.Content, "# Iterations") || strings.Contains(result.Content, "<!-- Managed by stew -->") {
		t.Fatalf("tail output included ledger preamble: %s", result.Content)
	}
}

func TestRunReturnsAllEntriesWhenLimitExceedsCount(t *testing.T) {
	content := ledgerContent(
		entry("2026-05-01T00:00:01Z", "First", "first body"),
		entry("2026-05-01T00:00:02Z", "Second", "second body"),
	)
	tmp := setupTailLedger(t, "iterations", content)

	result, err := Run(Options{
		TargetDir: tmp,
		Ledger:    "iterations",
		Limit:     10,
	})
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}

	want := entry("2026-05-01T00:00:01Z", "First", "first body") + "\n" +
		entry("2026-05-01T00:00:02Z", "Second", "second body")
	if result.Content != want {
		t.Fatalf("Content = %q, want %q", result.Content, want)
	}
}

func TestRunReturnsEmptyContentForLedgerWithoutEntries(t *testing.T) {
	tmp := setupTailLedger(t, "iterations", "# Iterations\n\n<!-- Managed by stew -->\n")

	result, err := Run(Options{
		TargetDir: tmp,
		Ledger:    "iterations",
		Limit:     5,
	})
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if result.Content != "" {
		t.Fatalf("Content = %q, want empty", result.Content)
	}
}

func TestRunRejectsInvalidLimits(t *testing.T) {
	tmp := setupTailLedger(t, "iterations", ledgerContent(entry("2026-05-01T00:00:01Z", "First", "body")))
	for _, limit := range []int{0, -1} {
		t.Run(fmt.Sprintf("%d", limit), func(t *testing.T) {
			_, err := Run(Options{
				TargetDir: tmp,
				Ledger:    "iterations",
				Limit:     limit,
			})
			if !errors.Is(err, ErrInvalidLimit) {
				t.Fatalf("error = %v, want ErrInvalidLimit", err)
			}
		})
	}
}

func TestRunIgnoresNonEntryH2HeadingsInsideEntries(t *testing.T) {
	first := entry("2026-05-01T00:00:01Z", "First", "first body\n\n## Notes\n\nnot an entry")
	second := entry("2026-05-01T00:00:02Z", "Second", "second body")
	third := entry("2026-05-01T00:00:03Z", "Third", "third body")
	tmp := setupTailLedger(t, "iterations", ledgerContent(first, second, third))

	result, err := Run(Options{
		TargetDir: tmp,
		Ledger:    "iterations",
		Limit:     2,
	})
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}

	want := second + "\n" + third
	if result.Content != want {
		t.Fatalf("Content = %q, want %q", result.Content, want)
	}
}

func TestRunSupportsDefaultTargetDir(t *testing.T) {
	content := ledgerContent(entry("2026-05-01T00:00:01Z", "First", "body"))
	tmp := setupTailLedger(t, "iterations", content)
	chdirForTailTest(t, tmp)

	result, err := Run(Options{
		Ledger: "iterations",
		Limit:  1,
	})
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if result.Content != entry("2026-05-01T00:00:01Z", "First", "body") {
		t.Fatalf("Content = %q", result.Content)
	}
}

func setupTailLedger(t *testing.T, ledger, content string) string {
	t.Helper()
	tmp := t.TempDir()
	stewDir := filepath.Join(tmp, ".stew")
	if err := os.MkdirAll(stewDir, 0o755); err != nil {
		t.Fatalf("mkdir .stew: %v", err)
	}
	if err := os.WriteFile(filepath.Join(stewDir, ledger+".spec.md"), []byte("# Spec\n"), 0o644); err != nil {
		t.Fatalf("write spec: %v", err)
	}
	if err := os.WriteFile(filepath.Join(stewDir, ledger+".md"), []byte(content), 0o644); err != nil {
		t.Fatalf("write ledger: %v", err)
	}
	return tmp
}

func ledgerContent(entries ...string) string {
	if len(entries) == 0 {
		return "# Iterations\n\n<!-- Managed by stew -->\n"
	}
	return "# Iterations\n\n<!-- Managed by stew -->\n\n" + strings.Join(entries, "\n")
}

func entry(timestamp, summary, body string) string {
	return "## " + timestamp + " — " + summary + "\n\n" +
		"**Prompt:** Test prompt\n\n" +
		body + "\n\n---\n"
}

func chdirForTailTest(t *testing.T, dir string) {
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
