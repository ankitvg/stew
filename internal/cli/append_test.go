package cli

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestAppendCommandAppendsMessageBody(t *testing.T) {
	tmp := setupCLILedger(t, "iterations")
	var out bytes.Buffer

	err := ExecuteWithIO([]string{
		"append", "iterations",
		"--path", tmp,
		"--prompt", "Implement append",
		"--summary", "Add append command",
		"-m", "CLI body",
	}, strings.NewReader(""), &out, &bytes.Buffer{})
	if err != nil {
		t.Fatalf("ExecuteWithIO() error = %v", err)
	}
	if got := out.String(); got != "Appended iterations\n" {
		t.Fatalf("stdout = %q", got)
	}

	content := readOnlyCLIEntry(t, tmp, "iterations")
	if !strings.Contains(content, "**Prompt:** Implement append\n\nCLI body\n\n---\n") {
		t.Fatalf("entry missing appended body: %s", content)
	}
}

func TestAppendCommandReadsDefaultStdinBody(t *testing.T) {
	tmp := setupCLILedger(t, "iterations")

	err := ExecuteWithIO([]string{
		"append", "iterations",
		"--path", tmp,
		"--prompt", "Use stdin",
		"--summary", "Read stdin",
	}, strings.NewReader("stdin body\n"), &bytes.Buffer{}, &bytes.Buffer{})
	if err != nil {
		t.Fatalf("ExecuteWithIO() error = %v", err)
	}

	content := readOnlyCLIEntry(t, tmp, "iterations")
	if !strings.Contains(content, "**Prompt:** Use stdin\n\nstdin body\n\n---\n") {
		t.Fatalf("entry missing stdin body: %s", content)
	}
}

func TestAppendCommandOutputsJSON(t *testing.T) {
	tmp := setupCLILedger(t, "iterations")
	var out bytes.Buffer

	err := ExecuteWithIO([]string{
		"append", "iterations",
		"--json",
		"--path", tmp,
		"--prompt", "Implement append",
		"--summary", "Add append command",
		"-m", "CLI body",
	}, strings.NewReader(""), &out, &bytes.Buffer{})
	if err != nil {
		t.Fatalf("ExecuteWithIO() error = %v", err)
	}

	var got map[string]string
	if err := json.Unmarshal(out.Bytes(), &got); err != nil {
		t.Fatalf("unmarshal JSON: %v\n%s", err, out.String())
	}
	if len(got) != 2 {
		t.Fatalf("JSON keys = %v, want only ledger and entryRef", keys(got))
	}
	if got["ledger"] != "iterations" {
		t.Fatalf("ledger = %q, want iterations", got["ledger"])
	}
	entries, err := os.ReadDir(filepath.Join(tmp, ".stew", "ledgers", "iterations"))
	if err != nil {
		t.Fatalf("read entry dir: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("entries len = %d, want 1", len(entries))
	}
	wantRef := "entry:iterations/" + entries[0].Name()
	if got["entryRef"] != wantRef {
		t.Fatalf("entryRef = %q, want %q", got["entryRef"], wantRef)
	}
	name := entries[0].Name()
	if !strings.HasSuffix(name, "-add-append-command.md") {
		t.Fatalf("entry filename did not include expected summary slug: %s", name)
	}
	prefix := strings.TrimSuffix(name, "-add-append-command.md")
	_, entryID, found := strings.Cut(prefix, "Z-")
	if !found || len(entryID) != 6 {
		t.Fatalf("entry filename did not include expected generated shape: %s", entries[0].Name())
	}
	assertNoPathFields(t, out.String(), tmp)
}

func TestAppendCommandRejectsMultipleExplicitBodySources(t *testing.T) {
	tmp := setupCLILedger(t, "iterations")
	bodyPath := filepath.Join(tmp, "body.md")
	if err := os.WriteFile(bodyPath, []byte("file body"), 0o644); err != nil {
		t.Fatalf("write body: %v", err)
	}

	err := ExecuteWithIO([]string{
		"append", "iterations",
		"--path", tmp,
		"--prompt", "Conflict",
		"--summary", "Reject conflict",
		"-m", "message body",
		"-F", bodyPath,
	}, strings.NewReader(""), &bytes.Buffer{}, &bytes.Buffer{})
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestAppendCommandRequiresPromptAndSummary(t *testing.T) {
	tmp := setupCLILedger(t, "iterations")

	err := ExecuteWithIO([]string{
		"append", "iterations",
		"--path", tmp,
		"-m", "body",
	}, strings.NewReader(""), &bytes.Buffer{}, &bytes.Buffer{})
	if err == nil {
		t.Fatalf("expected required flag error")
	}
	if !strings.Contains(err.Error(), "required flag") {
		t.Fatalf("error = %v, want required flag error", err)
	}
}

func setupCLILedger(t *testing.T, ledger string) string {
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

func readOnlyCLIEntry(t *testing.T, dir, ledger string) string {
	t.Helper()
	entryDir := filepath.Join(dir, ".stew", "ledgers", ledger)
	entries, err := os.ReadDir(entryDir)
	if err != nil {
		t.Fatalf("read entry dir: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("entries len = %d, want 1", len(entries))
	}
	return readCLIFile(t, filepath.Join(entryDir, entries[0].Name()))
}

func readCLIFile(t *testing.T, path string) string {
	t.Helper()
	bytes, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return string(bytes)
}
