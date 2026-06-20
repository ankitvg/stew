package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ankitvg/stew/internal/stewentry"
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

	ledgerInfo, err := os.Stat(filepath.Join(tmp, ".stew", "ledgers", "plans"))
	if err != nil {
		t.Fatalf("stat ledger storage: %v", err)
	}
	if !ledgerInfo.IsDir() {
		t.Fatalf("ledger storage should be a directory")
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

	want := ""
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

	want := ""
	if out.String() != want {
		t.Fatalf("stdout = %q, want %q", out.String(), want)
	}
}

func TestLedgerCatAllPrintsAllLedgers(t *testing.T) {
	tmp := setupCLIAllLedgers(t)
	zeta := cliEntry("2026-05-01T00:00:02Z", "Zeta", "zeta")
	alpha := cliEntry("2026-05-01T00:00:01Z", "Alpha", "alpha")
	writeCLIAllLedger(t, tmp, "zeta", zeta)
	writeCLIAllLedger(t, tmp, "alpha", alpha)
	var out bytes.Buffer

	err := ExecuteWithIO([]string{"ledger", "cat", "--all", "--path", tmp}, strings.NewReader(""), &out, &bytes.Buffer{})
	if err != nil {
		t.Fatalf("ExecuteWithIO() error = %v", err)
	}

	want := "# alpha\n\n" + alpha + "\n" +
		"# zeta\n\n" + zeta
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
		"For one ledger, the command writes concatenated entry markdown to stdout.",
		"--all, the command writes each ledger under a ledger-name section.",
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

func TestLedgerTailCommandOutputsJSON(t *testing.T) {
	tmp := setupCLILedger(t, "iterations")
	writeCLILedgerContent(t, tmp, "iterations", cliLedgerContent(
		cliEntry("2026-05-01T00:00:01Z", "First", "first body"),
		cliEntry("2026-05-01T00:00:02Z", "Second", "second body"),
		cliEntry("2026-05-01T00:00:03Z", "Third", "third body with <angle>"),
	))
	var out bytes.Buffer

	err := ExecuteWithIO([]string{"ledger", "tail", "iterations", "--json", "--path", tmp, "--limit", "2"}, strings.NewReader(""), &out, &bytes.Buffer{})
	if err != nil {
		t.Fatalf("ExecuteWithIO() error = %v", err)
	}

	ledger, entries := decodeTailJSON(t, out.String())
	if ledger != "iterations" {
		t.Fatalf("ledger = %q, want iterations", ledger)
	}
	assertTailEntryJSON(t, entries, []wantTailEntry{
		{
			timestamp: "2026-05-01T00:00:02Z",
			summary:   "Second",
			prompt:    "Test prompt",
			body:      "second body",
		},
		{
			timestamp: "2026-05-01T00:00:03Z",
			summary:   "Third",
			prompt:    "Test prompt",
			body:      "third body with <angle>",
		},
	})
	if strings.Contains(out.String(), `\u003c`) || strings.Contains(out.String(), `\u003e`) {
		t.Fatalf("JSON output should not HTML-escape angle brackets: %s", out.String())
	}
	assertNoPathFields(t, out.String(), tmp)
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

func TestLedgerTailAllOutputsJSON(t *testing.T) {
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

	err := ExecuteWithIO([]string{"ledger", "tail", "--all", "--json", "--path", tmp, "--limit", "2"}, strings.NewReader(""), &out, &bytes.Buffer{})
	if err != nil {
		t.Fatalf("ExecuteWithIO() error = %v", err)
	}

	ledgers := decodeTailAllJSON(t, out.String())
	if len(ledgers) != 2 {
		t.Fatalf("ledgers len = %d, want 2", len(ledgers))
	}
	if ledgers[0].ledger != "alpha" || ledgers[1].ledger != "zeta" {
		t.Fatalf("ledger order = %q, %q; want alpha, zeta", ledgers[0].ledger, ledgers[1].ledger)
	}
	assertTailEntryJSON(t, ledgers[0].entries, []wantTailEntry{
		{
			timestamp: "2026-05-01T00:00:02Z",
			summary:   "Alpha two",
			prompt:    "Test prompt",
			body:      "alpha two",
		},
		{
			timestamp: "2026-05-01T00:00:03Z",
			summary:   "Alpha three",
			prompt:    "Test prompt",
			body:      "alpha three",
		},
	})
	assertTailEntryJSON(t, ledgers[1].entries, []wantTailEntry{
		{
			timestamp: "2026-05-01T00:00:05Z",
			summary:   "Zeta two",
			prompt:    "Test prompt",
			body:      "zeta two",
		},
		{
			timestamp: "2026-05-01T00:00:06Z",
			summary:   "Zeta three",
			prompt:    "Test prompt",
			body:      "zeta three",
		},
	})
	assertNoPathFields(t, out.String(), tmp)
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
		"stew ledger tail iterations --json --limit 5",
		"stew ledger tail --all --limit 5",
		"stew ledger tail --all --json --limit 5",
		"--all",
		"--json",
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

type decodedTailLedger struct {
	ledger  string
	entries []map[string]string
}

type wantTailEntry struct {
	timestamp string
	summary   string
	prompt    string
	body      string
}

func decodeTailJSON(t *testing.T, output string) (string, []map[string]string) {
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
	if len(topLevel) != 2 {
		t.Fatalf("top-level keys = %v, want only ledger and entries", keys(topLevel))
	}
	rawLedger, ok := topLevel["ledger"]
	if !ok {
		t.Fatalf("top-level JSON missing ledger key: %v", keys(topLevel))
	}
	rawEntries, ok := topLevel["entries"]
	if !ok {
		t.Fatalf("top-level JSON missing entries key: %v", keys(topLevel))
	}

	var ledger string
	if err := json.Unmarshal(rawLedger, &ledger); err != nil {
		t.Fatalf("unmarshal ledger: %v", err)
	}
	var entries []map[string]string
	if err := json.Unmarshal(rawEntries, &entries); err != nil {
		t.Fatalf("unmarshal entries: %v", err)
	}
	return ledger, entries
}

func decodeTailAllJSON(t *testing.T, output string) []decodedTailLedger {
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

	var rawItems []map[string]json.RawMessage
	if err := json.Unmarshal(rawLedgers, &rawItems); err != nil {
		t.Fatalf("unmarshal ledgers: %v", err)
	}
	ledgers := make([]decodedTailLedger, 0, len(rawItems))
	for i, item := range rawItems {
		if len(item) != 2 {
			t.Fatalf("ledgers[%d] keys = %v, want only ledger and entries", i, keys(item))
		}
		var ledger string
		if err := json.Unmarshal(item["ledger"], &ledger); err != nil {
			t.Fatalf("unmarshal ledgers[%d].ledger: %v", i, err)
		}
		var entries []map[string]string
		if err := json.Unmarshal(item["entries"], &entries); err != nil {
			t.Fatalf("unmarshal ledgers[%d].entries: %v", i, err)
		}
		ledgers = append(ledgers, decodedTailLedger{ledger: ledger, entries: entries})
	}
	return ledgers
}

func assertTailEntryJSON(t *testing.T, got []map[string]string, want []wantTailEntry) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("entries len = %d, want %d", len(got), len(want))
	}
	for i := range want {
		if len(got[i]) != 4 {
			t.Fatalf("entries[%d] keys = %v, want only timestamp, summary, prompt, body", i, keys(got[i]))
		}
		if got[i]["timestamp"] != want[i].timestamp ||
			got[i]["summary"] != want[i].summary ||
			got[i]["prompt"] != want[i].prompt ||
			got[i]["body"] != want[i].body {
			t.Fatalf("entries[%d] = %#v, want timestamp=%q summary=%q prompt=%q body=%q", i, got[i], want[i].timestamp, want[i].summary, want[i].prompt, want[i].body)
		}
	}
}

func writeCLILedgerContent(t *testing.T, dir, ledger, content string) {
	t.Helper()
	entryDir := filepath.Join(dir, ".stew", "ledgers", ledger)
	if err := os.MkdirAll(entryDir, 0o755); err != nil {
		t.Fatalf("mkdir ledger storage: %v", err)
	}
	for _, entry := range stewentry.Parse(content) {
		name, err := stewentry.Filename(entry.Timestamp, entry.Summary)
		if err != nil {
			t.Fatalf("entry filename: %v", err)
		}
		if err := os.WriteFile(filepath.Join(entryDir, name), []byte(strings.TrimRight(entry.Content, "\n")+"\n"), 0o644); err != nil {
			t.Fatalf("write entry: %v", err)
		}
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
	writeCLILedgerContent(t, dir, ledger, content)
}

func cliLedgerContent(entries ...string) string {
	if len(entries) == 0 {
		return ""
	}
	return strings.Join(entries, "\n")
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
