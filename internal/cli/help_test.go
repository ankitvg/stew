package cli

import (
	"bytes"
	"strings"
	"testing"
)

func TestRootHelpIncludesAgentWorkflow(t *testing.T) {
	output := executeHelp(t, "--help")

	required := []string{
		"Stew maintains append-only markdown ledgers in .stew/.",
		`Run "stew help" to discover available commands.`,
		`Run "stew full-spec" to load the repository's ledger contract.`,
		`Run "stew <command> --help" before using a command you do not know.`,
		`Use "stew ledger new <name>" to create custom ledgers when needed.`,
		"stew ledger new plans --description",
		"stew append iterations --help",
	}
	assertHelpContains(t, output, required)
}

func TestLedgerHelpDocumentsCustomLedgers(t *testing.T) {
	output := executeHelp(t, "ledger", "--help")

	required := []string{
		"Manage Stew ledgers.",
		"Stew discovers ledgers from .stew/*.spec.md files.",
		`Use "stew ledger new" when`,
		"new",
	}
	assertHelpContains(t, output, required)
}

func TestLedgerNewHelpDocumentsCreationContract(t *testing.T) {
	output := executeHelp(t, "ledger", "new", "--help")

	required := []string{
		"Create a custom Stew ledger in an initialized repository.",
		"The command writes .stew/<name>.md and .stew/<name>.spec.md.",
		"Names must use lowercase ASCII letters, digits, and single hyphens only.",
		"omitted, the spec contains TODO guidance",
		"stew ledger new plans --description",
	}
	assertHelpContains(t, output, required)
}

func TestAppendHelpDocumentsBodySourcesAndExamples(t *testing.T) {
	output := executeHelp(t, "append", "--help")

	required := []string{
		"Append a new entry to .stew/<ledger>.md.",
		"The ledger must already be defined by .stew/<ledger>.spec.md.",
		"source: piped stdin, -m/--message, or -F/--file.",
		"printf 'Implemented the change and ran go test ./...'",
		"stew append iterations --prompt 'Small fix' --summary 'Record small fix' -m",
		"stew append decisions --prompt 'Choose storage model' --summary 'Use append-only ledgers' -F decision-entry.md",
	}
	assertHelpContains(t, output, required)
}

func TestInitHelpDocumentsManagedAgentsBlock(t *testing.T) {
	output := executeHelp(t, "init", "--help")

	required := []string{
		"creates the .stew/ directory, the default ledger/spec files, and",
		"the managed Stew block in AGENTS.md",
		`agents should use "stew help" for CLI discovery`,
		`Run "stew full-spec"`,
		"to load the repository's ledger contract",
	}
	assertHelpContains(t, output, required)
}

func TestFullSpecHelpDocumentsContractLoading(t *testing.T) {
	output := executeHelp(t, "full-spec", "--help")

	required := []string{
		"Print the full Stew contract for the target repository.",
		"The output starts with .stew/stew.spec.md",
		"Agents should run this before",
		"ledger's expected entry shape",
	}
	assertHelpContains(t, output, required)
}

func executeHelp(t *testing.T, args ...string) string {
	t.Helper()
	var out bytes.Buffer
	cmd := newRootCmd()
	cmd.SetOut(&out)
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetArgs(args)
	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute(%v) error = %v", args, err)
	}
	return out.String()
}

func assertHelpContains(t *testing.T, output string, required []string) {
	t.Helper()
	for _, want := range required {
		if !strings.Contains(output, want) {
			t.Fatalf("help output missing %q\n--- output ---\n%s", want, output)
		}
	}
}
