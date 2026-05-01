package cli

import (
	"bytes"
	"strings"
	"testing"
)

func TestInitCommandDoesNotExposeStewFilePaths(t *testing.T) {
	tmp := t.TempDir()
	var out bytes.Buffer

	err := ExecuteWithIO([]string{"init", "--path", tmp}, strings.NewReader(""), &out, &bytes.Buffer{})
	if err != nil {
		t.Fatalf("ExecuteWithIO() error = %v", err)
	}

	output := out.String()
	required := []string{
		"Initialized stew in ",
		"[created] Stew metadata",
		"[created] AGENTS.md",
	}
	assertHelpContains(t, output, required)
	if strings.Contains(output, ".stew") {
		t.Fatalf("init output should not expose stew file paths\n--- output ---\n%s", output)
	}
}
