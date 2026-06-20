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

func TestAppendCommandOutputsJSONWithLinks(t *testing.T) {
	tmp := setupCLILedger(t, "iterations")
	writeCLIRefTarget(t, filepath.Join(tmp, "internal", "cli", "ledger.go"), "package cli\n")
	writeCLIRefTarget(t, filepath.Join(tmp, "internal", "stewappend", "stewappend.go"), "package stewappend\n")
	var out bytes.Buffer

	err := ExecuteWithIO([]string{
		"append", "iterations",
		"--json",
		"--path", tmp,
		"--prompt", "Implement append",
		"--summary", "Add append command",
		"-m", "CLI body",
		"--link-file", `./internal\cli/ledger.go`,
		"--link-file", "internal/stewappend/stewappend.go",
	}, strings.NewReader(""), &out, &bytes.Buffer{})
	if err != nil {
		t.Fatalf("ExecuteWithIO() error = %v", err)
	}

	var topLevel map[string]json.RawMessage
	if err := json.Unmarshal(out.Bytes(), &topLevel); err != nil {
		t.Fatalf("unmarshal JSON: %v\n%s", err, out.String())
	}
	if len(topLevel) != 3 {
		t.Fatalf("JSON keys = %v, want ledger, entryRef, links", keys(topLevel))
	}
	var entryRef string
	if err := json.Unmarshal(topLevel["entryRef"], &entryRef); err != nil {
		t.Fatalf("unmarshal entryRef: %v", err)
	}
	var links []map[string]string
	if err := json.Unmarshal(topLevel["links"], &links); err != nil {
		t.Fatalf("unmarshal links: %v", err)
	}
	if len(links) != 2 {
		t.Fatalf("links len = %d, want 2", len(links))
	}
	wantTargets := []string{"file:internal/cli/ledger.go", "file:internal/stewappend/stewappend.go"}
	for i, wantTarget := range wantTargets {
		if len(links[i]) != 3 {
			t.Fatalf("links[%d] keys = %v, want source, target, createdAt", i, keys(links[i]))
		}
		if links[i]["source"] != entryRef || links[i]["target"] != wantTarget || links[i]["createdAt"] == "" {
			t.Fatalf("links[%d] = %#v", i, links[i])
		}
	}
	linkFiles, err := os.ReadDir(filepath.Join(tmp, ".stew", "links"))
	if err != nil {
		t.Fatalf("read links: %v", err)
	}
	if len(linkFiles) != 2 {
		t.Fatalf("link files len = %d, want 2", len(linkFiles))
	}
	assertNoPathFields(t, out.String(), tmp)
}

func TestAppendCommandLinksFilesWithoutChangingDefaultOutput(t *testing.T) {
	tmp := setupCLILedger(t, "iterations")
	bodyPath := filepath.Join(tmp, "body.md")
	if err := os.WriteFile(bodyPath, []byte("file body"), 0o644); err != nil {
		t.Fatalf("write body: %v", err)
	}
	writeCLIRefTarget(t, filepath.Join(tmp, "internal", "cli", "ledger.go"), "package cli\n")
	var out bytes.Buffer

	err := ExecuteWithIO([]string{
		"append", "iterations",
		"--path", tmp,
		"--prompt", "Implement append",
		"--summary", "Add append command",
		"-F", bodyPath,
		"--link-file", "internal/cli/ledger.go",
	}, strings.NewReader(""), &out, &bytes.Buffer{})
	if err != nil {
		t.Fatalf("ExecuteWithIO() error = %v", err)
	}
	if out.String() != "Appended iterations\n" {
		t.Fatalf("stdout = %q", out.String())
	}
	content := readOnlyCLIEntry(t, tmp, "iterations")
	if !strings.Contains(content, "\n\nfile body\n\n---\n") {
		t.Fatalf("entry missing file body: %s", content)
	}
	linkFiles, err := os.ReadDir(filepath.Join(tmp, ".stew", "links"))
	if err != nil {
		t.Fatalf("read links: %v", err)
	}
	if len(linkFiles) != 1 {
		t.Fatalf("link files len = %d, want 1", len(linkFiles))
	}
}

func TestAppendCommandRejectsMissingLinkedFileBeforeEntryCreation(t *testing.T) {
	tmp := setupCLILedger(t, "iterations")

	err := ExecuteWithIO([]string{
		"append", "iterations",
		"--path", tmp,
		"--prompt", "Implement append",
		"--summary", "Add append command",
		"-m", "CLI body",
		"--link-file", "internal/cli/missing.go",
	}, strings.NewReader(""), &bytes.Buffer{}, &bytes.Buffer{})
	if err == nil {
		t.Fatalf("expected missing linked file error")
	}
	entries, readErr := os.ReadDir(filepath.Join(tmp, ".stew", "ledgers", "iterations"))
	if readErr != nil {
		t.Fatalf("read entries: %v", readErr)
	}
	if len(entries) != 0 {
		t.Fatalf("entries len = %d, want 0", len(entries))
	}
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
