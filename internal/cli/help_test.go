package cli

import (
	"bytes"
	"strings"
	"testing"
)

func TestRootHelpIncludesAgentWorkflow(t *testing.T) {
	output := executeHelp(t, "--help")

	required := []string{
		"Stew maintains append-only markdown ledger entries through the stew CLI.",
		`Run "stew help" to discover available commands.`,
		`Run "stew full-spec" to load the repository's ledger contract.`,
		`Run "stew <command> --help" before using a command you do not know.`,
		`Use "stew ledger new <name>" to create custom ledgers when needed.`,
		"stew ledger new plans --description",
		"stew append iterations --help",
	}
	assertHelpContains(t, output, required)
	if strings.Contains(output, "completion") {
		t.Fatalf("root help should not include Cobra completion command\n--- output ---\n%s", output)
	}
}

func TestLedgerHelpDocumentsCustomLedgers(t *testing.T) {
	output := executeHelp(t, "ledger", "--help")

	required := []string{
		"Manage Stew ledgers.",
		"Stew discovers ledgers from ledger specs.",
		`Use "stew ledger new" when`,
		"new",
		"cat",
		"Print one or all stew ledgers",
	}
	assertHelpContains(t, output, required)
}

func TestLedgerNewHelpDocumentsCreationContract(t *testing.T) {
	output := executeHelp(t, "ledger", "new", "--help")

	required := []string{
		"Create a custom Stew ledger in an initialized repository.",
		"The command creates an atomic entry storage directory and its matching ledger",
		"Names must use lowercase ASCII letters, digits, and single hyphens only.",
		"omitted, the spec contains TODO guidance",
		"stew ledger new plans --description",
	}
	assertHelpContains(t, output, required)
	if strings.Contains(output, "--path") {
		t.Fatalf("ledger new help should not expose --path\n--- output ---\n%s", output)
	}
}

func TestMigrateHelpDocumentsCommands(t *testing.T) {
	output := executeHelp(t, "migrate", "--help")

	required := []string{
		"Migrate Stew metadata between storage formats.",
		"atomic-entries",
		"Split monolithic ledger files into atomic entry files",
	}
	assertHelpContains(t, output, required)
}

func TestAppendHelpDocumentsBodySourcesAndExamples(t *testing.T) {
	output := executeHelp(t, "append", "--help")

	required := []string{
		"Append a new entry to a Stew ledger.",
		"The ledger must already be defined by a ledger spec.",
		"source: piped stdin, -m/--message, or -F/--file.",
		"Use --json when an agent or script needs the created entry ref.",
		"printf 'Implemented the change and ran go test ./...'",
		"stew append iterations --prompt 'Small fix' --summary 'Record small fix' -m",
		"stew append decisions --prompt 'Choose storage model' --summary 'Use append-only ledgers' -F decision-entry.md",
		"--json",
	}
	assertHelpContains(t, output, required)
}

func TestInitHelpDocumentsManagedAgentsBlock(t *testing.T) {
	output := executeHelp(t, "init", "--help")

	required := []string{
		"creates Stew metadata, the default ledger storage/specs, and",
		"Stew block in AGENTS.md",
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
		"The output starts with the base Stew spec",
		"Agents should run this before",
		"expected entry shape",
	}
	assertHelpContains(t, output, required)
}

func executeHelp(t *testing.T, args ...string) string {
	t.Helper()
	var out bytes.Buffer
	if err := ExecuteWithIO(args, strings.NewReader(""), &out, &bytes.Buffer{}); err != nil {
		t.Fatalf("ExecuteWithIO(%v) error = %v", args, err)
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
