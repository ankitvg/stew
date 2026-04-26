package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestAppendCommandAppendsMessageBody(t *testing.T) {
	tmp := setupCLILedger(t, "iterations")
	var out bytes.Buffer

	cmd := newRootCmd()
	cmd.SetOut(&out)
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetArgs([]string{
		"append", "iterations",
		"--path", tmp,
		"--prompt", "Implement append",
		"--summary", "Add append command",
		"-m", "CLI body",
	})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if got := out.String(); got != "Appended "+filepath.Join(".stew", "iterations.md")+"\n" {
		t.Fatalf("stdout = %q", got)
	}

	content := readCLIFile(t, filepath.Join(tmp, ".stew", "iterations.md"))
	if !strings.Contains(content, "**Prompt:** Implement append\n\nCLI body\n\n---\n") {
		t.Fatalf("ledger missing appended body: %s", content)
	}
}

func TestAppendCommandReadsDefaultStdinBody(t *testing.T) {
	tmp := setupCLILedger(t, "iterations")

	cmd := newRootCmd()
	cmd.SetOut(&bytes.Buffer{})
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetIn(strings.NewReader("stdin body\n"))
	cmd.SetArgs([]string{
		"append", "iterations",
		"--path", tmp,
		"--prompt", "Use stdin",
		"--summary", "Read stdin",
	})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	content := readCLIFile(t, filepath.Join(tmp, ".stew", "iterations.md"))
	if !strings.Contains(content, "**Prompt:** Use stdin\n\nstdin body\n\n---\n") {
		t.Fatalf("ledger missing stdin body: %s", content)
	}
}

func TestAppendCommandRejectsMultipleExplicitBodySources(t *testing.T) {
	tmp := setupCLILedger(t, "iterations")
	bodyPath := filepath.Join(tmp, "body.md")
	if err := os.WriteFile(bodyPath, []byte("file body"), 0o644); err != nil {
		t.Fatalf("write body: %v", err)
	}

	cmd := newRootCmd()
	cmd.SetOut(&bytes.Buffer{})
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetArgs([]string{
		"append", "iterations",
		"--path", tmp,
		"--prompt", "Conflict",
		"--summary", "Reject conflict",
		"-m", "message body",
		"-F", bodyPath,
	})

	if err := cmd.Execute(); err == nil {
		t.Fatalf("expected error")
	}
}

func TestAppendCommandRequiresPromptAndSummary(t *testing.T) {
	tmp := setupCLILedger(t, "iterations")

	cmd := newRootCmd()
	cmd.SetOut(&bytes.Buffer{})
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetArgs([]string{
		"append", "iterations",
		"--path", tmp,
		"-m", "body",
	})

	err := cmd.Execute()
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
	if err := os.WriteFile(filepath.Join(stewDir, ledger+".md"), []byte("# Iterations\n\n<!-- Managed by stew -->\n"), 0o644); err != nil {
		t.Fatalf("write ledger: %v", err)
	}
	return tmp
}

func readCLIFile(t *testing.T, path string) string {
	t.Helper()
	bytes, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return string(bytes)
}
