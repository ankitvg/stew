package cli

import (
	"bytes"
	"fmt"
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
	wantOut := "Created ledger plans\n"
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
	if !strings.Contains(fullSpecOut.String(), "<!-- plans spec -->") {
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

func TestLedgerCatAllPrintsAllLedgers(t *testing.T) {
	tmp := setupCLIAllLedgers(t)
	writeCLIAllLedger(t, tmp, "zeta", "# Zeta\n\n<!-- Managed by stew -->\n")
	writeCLIAllLedger(t, tmp, "alpha", "# Alpha\n\n<!-- Managed by stew -->\n")
	var out bytes.Buffer

	err := ExecuteWithIO([]string{"ledger", "cat", "--all", "--path", tmp}, strings.NewReader(""), &out, &bytes.Buffer{})
	if err != nil {
		t.Fatalf("ExecuteWithIO() error = %v", err)
	}

	want := "# alpha\n\n# Alpha\n\n<!-- Managed by stew -->\n\n" +
		"# zeta\n\n# Zeta\n\n<!-- Managed by stew -->\n"
	if out.String() != want {
		t.Fatalf("stdout = %q, want %q", out.String(), want)
	}
	if strings.Contains(out.String(), ".stew") {
		t.Fatalf("cat --all output exposed ledger path: %s", out.String())
	}
}

func TestLedgerCatHelpDocumentsRawOutput(t *testing.T) {
	output := executeHelp(t, "ledger", "cat", "--help")

	required := []string{
		"Print Stew ledger content.",
		"For one ledger, the command writes raw ledger markdown to stdout.",
		"With --all",
		"stew ledger cat --all",
		"stew ledger cat iterations | grep",
		"--all",
		"--path string",
	}
	assertHelpContains(t, output, required)

	helpOutput := executeHelp(t, "help", "ledger", "cat")
	assertHelpContains(t, helpOutput, required)
}

func TestLedgerTailCommandPrintsRecentEntries(t *testing.T) {
	tmp := setupCLILedger(t, "iterations")
	writeCLILedgerContent(t, tmp, "iterations", cliLedgerContent(
		cliEntry("2026-05-01T00:00:01Z", "First", "first body"),
		cliEntry("2026-05-01T00:00:02Z", "Second", "second body"),
		cliEntry("2026-05-01T00:00:03Z", "Third", "third body"),
	))
	var out bytes.Buffer

	err := ExecuteWithIO([]string{"ledger", "tail", "iterations", "--path", tmp, "--limit", "2"}, strings.NewReader(""), &out, &bytes.Buffer{})
	if err != nil {
		t.Fatalf("ExecuteWithIO() error = %v", err)
	}

	want := cliEntry("2026-05-01T00:00:02Z", "Second", "second body") + "\n" +
		cliEntry("2026-05-01T00:00:03Z", "Third", "third body")
	if out.String() != want {
		t.Fatalf("stdout = %q, want %q", out.String(), want)
	}
	if strings.Contains(out.String(), "# Iterations") || strings.Contains(out.String(), "<!-- Managed by stew -->") {
		t.Fatalf("tail output included ledger preamble: %s", out.String())
	}
}

func TestLedgerTailCommandUsesDefaultLimit(t *testing.T) {
	tmp := setupCLILedger(t, "iterations")
	entries := make([]string, 0, 12)
	for i := 1; i <= 12; i++ {
		entries = append(entries, cliEntry(fmt.Sprintf("2026-05-01T00:00:%02dZ", i), fmt.Sprintf("Entry %02d", i), fmt.Sprintf("body %02d", i)))
	}
	writeCLILedgerContent(t, tmp, "iterations", cliLedgerContent(entries...))
	var out bytes.Buffer

	err := ExecuteWithIO([]string{"ledger", "tail", "iterations", "--path", tmp}, strings.NewReader(""), &out, &bytes.Buffer{})
	if err != nil {
		t.Fatalf("ExecuteWithIO() error = %v", err)
	}

	if strings.Contains(out.String(), "Entry 01") || strings.Contains(out.String(), "Entry 02") {
		t.Fatalf("default tail included entries beyond last 10: %s", out.String())
	}
	if !strings.Contains(out.String(), "Entry 03") || !strings.Contains(out.String(), "Entry 12") {
		t.Fatalf("default tail missing expected last 10 entries: %s", out.String())
	}
}

func TestLedgerTailAllPrintsRecentEntriesPerLedger(t *testing.T) {
	tmp := setupCLIAllLedgers(t)
	writeCLIAllLedger(t, tmp, "zeta", cliLedgerContent(
		cliEntry("2026-05-01T00:00:04Z", "Zeta one", "zeta one"),
		cliEntry("2026-05-01T00:00:05Z", "Zeta two", "zeta two"),
		cliEntry("2026-05-01T00:00:06Z", "Zeta three", "zeta three"),
	))
	writeCLIAllLedger(t, tmp, "alpha", cliLedgerContent(
		cliEntry("2026-05-01T00:00:01Z", "Alpha one", "alpha one"),
		cliEntry("2026-05-01T00:00:02Z", "Alpha two", "alpha two"),
		cliEntry("2026-05-01T00:00:03Z", "Alpha three", "alpha three"),
	))
	var out bytes.Buffer

	err := ExecuteWithIO([]string{"ledger", "tail", "--all", "--path", tmp, "--limit", "2"}, strings.NewReader(""), &out, &bytes.Buffer{})
	if err != nil {
		t.Fatalf("ExecuteWithIO() error = %v", err)
	}

	output := out.String()
	if !strings.HasPrefix(output, "# alpha\n\n") {
		t.Fatalf("tail --all output did not start with alpha section: %s", output)
	}
	if strings.Index(output, "# alpha") > strings.Index(output, "# zeta") {
		t.Fatalf("tail --all output was not sorted by ledger name: %s", output)
	}
	if strings.Contains(output, "Alpha one") || strings.Contains(output, "Zeta one") {
		t.Fatalf("tail --all included entries outside per-ledger limit: %s", output)
	}
	for _, want := range []string{"Alpha two", "Alpha three", "Zeta two", "Zeta three"} {
		if !strings.Contains(output, want) {
			t.Fatalf("tail --all missing %q: %s", want, output)
		}
	}
	if strings.Contains(output, "# Iterations") || strings.Contains(output, "<!-- Managed by stew -->") {
		t.Fatalf("tail --all output included ledger preamble: %s", output)
	}
	if strings.Contains(output, ".stew") {
		t.Fatalf("tail --all output exposed ledger path: %s", output)
	}
}

func TestLedgerTailCommandRejectsInvalidLimit(t *testing.T) {
	tmp := setupCLILedger(t, "iterations")

	err := ExecuteWithIO([]string{"ledger", "tail", "iterations", "--path", tmp, "--limit", "0"}, strings.NewReader(""), &bytes.Buffer{}, &bytes.Buffer{})
	if err == nil {
		t.Fatalf("expected error")
	}
	if !strings.Contains(err.Error(), "invalid limit") {
		t.Fatalf("error = %v, want invalid limit", err)
	}
}

func TestLedgerReadAllRejectsLedgerArgument(t *testing.T) {
	tmp := setupCLILedger(t, "iterations")

	for _, args := range [][]string{
		{"ledger", "cat", "iterations", "--all", "--path", tmp},
		{"ledger", "tail", "iterations", "--all", "--path", tmp},
	} {
		t.Run(strings.Join(args[:3], " "), func(t *testing.T) {
			err := ExecuteWithIO(args, strings.NewReader(""), &bytes.Buffer{}, &bytes.Buffer{})
			if err == nil {
				t.Fatalf("expected error")
			}
			if !strings.Contains(err.Error(), "cannot use --all with a ledger argument") {
				t.Fatalf("error = %v, want mutual exclusion error", err)
			}
		})
	}
}

func TestLedgerReadCommandsStillRequireLedgerWithoutAll(t *testing.T) {
	for _, args := range [][]string{
		{"ledger", "cat"},
		{"ledger", "tail"},
	} {
		t.Run(strings.Join(args, " "), func(t *testing.T) {
			err := ExecuteWithIO(args, strings.NewReader(""), &bytes.Buffer{}, &bytes.Buffer{})
			if err == nil {
				t.Fatalf("expected error")
			}
			if !strings.Contains(err.Error(), "accepts 1 arg(s), received 0") {
				t.Fatalf("error = %v, want missing ledger error", err)
			}
		})
	}
}

func TestLedgerTailHelpDocumentsEntryAwareOutput(t *testing.T) {
	output := executeHelp(t, "ledger", "tail", "--help")

	required := []string{
		"Print recent Stew ledger entries.",
		"For one ledger, the command writes the last N entries to stdout.",
		"With --all",
		"stew ledger tail iterations --limit 5",
		"stew ledger tail --all --limit 5",
		"--all",
		"--limit int",
		"--path string",
	}
	assertHelpContains(t, output, required)

	helpOutput := executeHelp(t, "help", "ledger", "tail")
	assertHelpContains(t, helpOutput, required)
}

func TestLedgerHelpIncludesReadCommands(t *testing.T) {
	output := executeHelp(t, "ledger", "--help")

	required := []string{
		"cat",
		"Print one or all stew ledgers",
		"new",
		"Create a custom stew ledger",
		"tail",
		"Print recent ledger entries",
	}
	assertHelpContains(t, output, required)
}

func writeCLILedgerContent(t *testing.T, dir, ledger, content string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(dir, ".stew", ledger+".md"), []byte(content), 0o644); err != nil {
		t.Fatalf("write ledger: %v", err)
	}
}

func setupCLIAllLedgers(t *testing.T) string {
	t.Helper()
	return setupCLIStewDir(t)
}

func writeCLIAllLedger(t *testing.T, dir, ledger, content string) {
	t.Helper()
	stewDir := filepath.Join(dir, ".stew")
	if err := os.WriteFile(filepath.Join(stewDir, ledger+".spec.md"), []byte("# "+ledger+" Spec\n\nDescription.\n"), 0o644); err != nil {
		t.Fatalf("write spec: %v", err)
	}
	if err := os.WriteFile(filepath.Join(stewDir, ledger+".md"), []byte(content), 0o644); err != nil {
		t.Fatalf("write ledger: %v", err)
	}
}

func cliLedgerContent(entries ...string) string {
	if len(entries) == 0 {
		return "# Iterations\n\n<!-- Managed by stew -->\n"
	}
	return "# Iterations\n\n<!-- Managed by stew -->\n\n" + strings.Join(entries, "\n")
}

func cliEntry(timestamp, summary, body string) string {
	return "## " + timestamp + " — " + summary + "\n\n" +
		"**Prompt:** Test prompt\n\n" +
		body + "\n\n---\n"
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
