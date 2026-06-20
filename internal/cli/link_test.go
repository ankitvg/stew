package cli

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/ankitvg/stew/internal/stewlink"
	"github.com/ankitvg/stew/internal/stewref"
)

func TestLinkListCommandPrintsLinksForSource(t *testing.T) {
	tmp, source, target := setupCLILink(t)
	var out bytes.Buffer

	err := ExecuteWithIO([]string{"link", "list", source.String(), "--path", tmp}, strings.NewReader(""), &out, &bytes.Buffer{})
	if err != nil {
		t.Fatalf("ExecuteWithIO() error = %v", err)
	}

	want := source.String() + " -> " + target.String() + "\n"
	if out.String() != want {
		t.Fatalf("stdout = %q, want %q", out.String(), want)
	}
}

func TestLinkListCommandPrintsLinksForTarget(t *testing.T) {
	tmp, source, target := setupCLILink(t)
	var out bytes.Buffer

	err := ExecuteWithIO([]string{"link", "list", target.String(), "--path", tmp}, strings.NewReader(""), &out, &bytes.Buffer{})
	if err != nil {
		t.Fatalf("ExecuteWithIO() error = %v", err)
	}

	want := source.String() + " -> " + target.String() + "\n"
	if out.String() != want {
		t.Fatalf("stdout = %q, want %q", out.String(), want)
	}
}

func TestLinkListCommandOutputsJSON(t *testing.T) {
	tmp, source, target := setupCLILink(t)
	var out bytes.Buffer

	err := ExecuteWithIO([]string{"link", "list", target.String(), "--json", "--path", tmp}, strings.NewReader(""), &out, &bytes.Buffer{})
	if err != nil {
		t.Fatalf("ExecuteWithIO() error = %v", err)
	}

	var topLevel map[string]json.RawMessage
	if err := json.Unmarshal(out.Bytes(), &topLevel); err != nil {
		t.Fatalf("unmarshal JSON: %v\n%s", err, out.String())
	}
	if len(topLevel) != 2 {
		t.Fatalf("top-level keys = %v, want only ref and links", keys(topLevel))
	}
	var gotRef string
	if err := json.Unmarshal(topLevel["ref"], &gotRef); err != nil {
		t.Fatalf("unmarshal ref: %v", err)
	}
	if gotRef != target.String() {
		t.Fatalf("ref = %q, want %q", gotRef, target.String())
	}
	var links []map[string]string
	if err := json.Unmarshal(topLevel["links"], &links); err != nil {
		t.Fatalf("unmarshal links: %v", err)
	}
	if len(links) != 1 {
		t.Fatalf("links len = %d, want 1", len(links))
	}
	if len(links[0]) != 3 {
		t.Fatalf("link keys = %v, want only source, target, createdAt", keys(links[0]))
	}
	if links[0]["source"] != source.String() || links[0]["target"] != target.String() || links[0]["createdAt"] != "2026-05-01T12:05:07Z" {
		t.Fatalf("link = %#v", links[0])
	}
	assertNoPathFields(t, out.String(), tmp)
}

func TestLinkHelpDocumentsListCommand(t *testing.T) {
	output := executeHelp(t, "link", "--help")
	assertHelpContains(t, output, []string{
		"Manage Stew links.",
		"Links are append-only relationships between refs.",
		"Use append --link-file to create links.",
		"list",
		"List links for a ref",
	})

	listOutput := executeHelp(t, "link", "list", "--help")
	assertHelpContains(t, listOutput, []string{
		"List Stew links for a ref.",
		"source -> target",
		"stew link list file:internal/cli/ledger.go",
		"--json",
		"--path string",
	})
}

func setupCLILink(t *testing.T) (string, stewref.Ref, stewref.Ref) {
	t.Helper()
	tmp := setupCLILedger(t, "iterations")
	writeCLILedgerContent(t, tmp, "iterations", cliLedgerContent(
		cliEntry("2026-05-01T00:00:01Z", "First", "body"),
	))
	writeCLIRefTarget(t, filepath.Join(tmp, "internal", "cli", "ledger.go"), "package cli\n")

	source, err := stewref.Entry("iterations", "2026-05-01T000001Z-first.md")
	if err != nil {
		t.Fatalf("entry ref: %v", err)
	}
	target, err := stewref.File("internal/cli/ledger.go")
	if err != nil {
		t.Fatalf("file ref: %v", err)
	}
	_, err = stewlink.Create(stewlink.CreateOptions{
		TargetDir: tmp,
		Source:    source,
		Target:    target,
		Now: func() time.Time {
			return time.Date(2026, 5, 1, 12, 5, 7, 0, time.UTC)
		},
		NewID: func() (string, error) {
			return "k7p3qx", nil
		},
	})
	if err != nil {
		t.Fatalf("create link: %v", err)
	}
	return tmp, source, target
}

func writeCLIRefTarget(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("mkdir target: %v", err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write target: %v", err)
	}
}
