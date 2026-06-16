package stewledgerall

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ankitvg/stew/internal/stewentry"
	"github.com/ankitvg/stew/internal/stewledgercat"
)

func TestCatUsesDiscoveredLedgerOrder(t *testing.T) {
	tmp := setupAllStewDir(t)
	zeta := allEntry("2026-05-01T00:00:02Z", "Zeta")
	alpha := allEntry("2026-05-01T00:00:01Z", "Alpha")
	writeAllLedger(t, tmp, "zeta", zeta)
	writeAllLedger(t, tmp, "alpha", alpha)

	result, err := Cat(Options{TargetDir: tmp})
	if err != nil {
		t.Fatalf("Cat() error = %v", err)
	}

	if got, want := sectionNames(result.Sections), []string{"alpha", "zeta"}; strings.Join(got, ",") != strings.Join(want, ",") {
		t.Fatalf("section names = %#v, want %#v", got, want)
	}
	if result.Sections[0].Content != alpha || result.Sections[1].Content != zeta {
		t.Fatalf("sections = %#v", result.Sections)
	}
}

func TestTailAppliesLimitPerLedger(t *testing.T) {
	tmp := setupAllStewDir(t)
	writeAllLedger(t, tmp, "alpha", allLedgerContent(
		allEntry("2026-05-01T00:00:01Z", "Alpha one"),
		allEntry("2026-05-01T00:00:02Z", "Alpha two"),
		allEntry("2026-05-01T00:00:03Z", "Alpha three"),
	))
	writeAllLedger(t, tmp, "zeta", allLedgerContent(
		allEntry("2026-05-01T00:00:04Z", "Zeta one"),
		allEntry("2026-05-01T00:00:05Z", "Zeta two"),
		allEntry("2026-05-01T00:00:06Z", "Zeta three"),
	))

	result, err := Tail(Options{TargetDir: tmp, Limit: 2})
	if err != nil {
		t.Fatalf("Tail() error = %v", err)
	}

	if got, want := sectionNames(result.Sections), []string{"alpha", "zeta"}; strings.Join(got, ",") != strings.Join(want, ",") {
		t.Fatalf("section names = %#v, want %#v", got, want)
	}
	alpha := result.Sections[0].Content
	if strings.Contains(alpha, "Alpha one") || !strings.Contains(alpha, "Alpha two") || !strings.Contains(alpha, "Alpha three") {
		t.Fatalf("alpha section did not apply per-ledger limit: %s", alpha)
	}
	zeta := result.Sections[1].Content
	if strings.Contains(zeta, "Zeta one") || !strings.Contains(zeta, "Zeta two") || !strings.Contains(zeta, "Zeta three") {
		t.Fatalf("zeta section did not apply per-ledger limit: %s", zeta)
	}
}

func TestTailIncludesEmptyKnownLedgerSection(t *testing.T) {
	tmp := setupAllStewDir(t)
	writeAllLedger(t, tmp, "alpha", "")

	result, err := Tail(Options{TargetDir: tmp, Limit: 5})
	if err != nil {
		t.Fatalf("Tail() error = %v", err)
	}

	if len(result.Sections) != 1 {
		t.Fatalf("sections len = %d, want 1", len(result.Sections))
	}
	if result.Sections[0].Name != "alpha" || result.Sections[0].Content != "" {
		t.Fatalf("section = %#v, want empty alpha section", result.Sections[0])
	}
}

func TestCatFailsWhenDiscoveredLedgerContentIsMissing(t *testing.T) {
	tmp := setupAllStewDir(t)
	writeAllSpec(t, tmp, "alpha")

	_, err := Cat(Options{TargetDir: tmp})
	if !errors.Is(err, stewledgercat.ErrMissingLedger) {
		t.Fatalf("error = %v, want ErrMissingLedger", err)
	}
	if !strings.Contains(err.Error(), `read ledger "alpha"`) {
		t.Fatalf("error = %v, want ledger context", err)
	}
}

func TestTailFailsWhenDiscoveredLedgerContentIsMissing(t *testing.T) {
	tmp := setupAllStewDir(t)
	writeAllSpec(t, tmp, "alpha")

	_, err := Tail(Options{TargetDir: tmp, Limit: 5})
	if !errors.Is(err, stewledgercat.ErrMissingLedger) {
		t.Fatalf("error = %v, want ErrMissingLedger", err)
	}
	if !strings.Contains(err.Error(), `read ledger "alpha"`) {
		t.Fatalf("error = %v, want ledger context", err)
	}
}

func TestTailRejectsInvalidLimit(t *testing.T) {
	tmp := setupAllStewDir(t)
	writeAllLedger(t, tmp, "alpha", "")

	_, err := Tail(Options{TargetDir: tmp, Limit: 0})
	if err == nil || !strings.Contains(err.Error(), "invalid limit") {
		t.Fatalf("error = %v, want invalid limit", err)
	}
}

func setupAllStewDir(t *testing.T) string {
	t.Helper()
	tmp := t.TempDir()
	if err := os.Mkdir(filepath.Join(tmp, ".stew"), 0o755); err != nil {
		t.Fatalf("mkdir .stew: %v", err)
	}
	return tmp
}

func writeAllLedger(t *testing.T, dir, ledger, content string) {
	t.Helper()
	writeAllSpec(t, dir, ledger)
	entryDir := filepath.Join(dir, ".stew", "ledgers", ledger)
	if err := os.MkdirAll(entryDir, 0o755); err != nil {
		t.Fatalf("mkdir ledger %s: %v", ledger, err)
	}
	for _, entry := range stewentry.Parse(content) {
		name, err := stewentry.Filename(entry.Timestamp, entry.Summary)
		if err != nil {
			t.Fatalf("entry filename: %v", err)
		}
		if err := os.WriteFile(filepath.Join(entryDir, name), []byte(strings.TrimRight(entry.Content, "\n")+"\n"), 0o644); err != nil {
			t.Fatalf("write entry %s: %v", ledger, err)
		}
	}
}

func writeAllSpec(t *testing.T, dir, ledger string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(dir, ".stew", ledger+".spec.md"), []byte("# "+ledger+" Spec\n\nDescription.\n"), 0o644); err != nil {
		t.Fatalf("write spec %s: %v", ledger, err)
	}
}

func sectionNames(sections []Section) []string {
	names := make([]string, 0, len(sections))
	for _, section := range sections {
		names = append(names, section.Name)
	}
	return names
}

func allLedgerContent(entries ...string) string {
	return "# Ledger\n\n<!-- Managed by stew -->\n\n" + strings.Join(entries, "\n")
}

func allEntry(timestamp, summary string) string {
	return "## " + timestamp + " — " + summary + "\n\n" +
		"**Prompt:** Test prompt\n\n" +
		"Body\n\n---\n"
}
