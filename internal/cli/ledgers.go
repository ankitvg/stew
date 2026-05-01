package cli

import (
	"encoding/json"
	"fmt"
	"text/tabwriter"

	"github.com/ankitvg/stew/internal/stewledgers"
)

func runLedgers(ctx cliContext, args []string) error {
	if wantsHelp(args) {
		fmt.Fprint(ctx.out, ledgersHelp)
		return nil
	}

	var targetPath string
	var jsonOutput bool

	flags := newFlagSet("ledgers")
	flags.BoolVar(&jsonOutput, "json", false, "Print JSON output")
	flags.StringVar(&targetPath, "path", ".", "Target directory")

	positionals, err := parseInterspersedFlags(flags, args, map[string]flagKind{
		"json": boolFlag,
		"path": argFlag,
	})
	if err != nil {
		return err
	}
	if err := exactArgs(positionals, 0); err != nil {
		return err
	}

	result, err := stewledgers.List(stewledgers.Options{TargetDir: targetPath})
	if err != nil {
		return err
	}

	if jsonOutput {
		return json.NewEncoder(ctx.out).Encode(ledgersJSONResponseFrom(result))
	}

	writer := tabwriter.NewWriter(ctx.out, 0, 0, 2, ' ', 0)
	for _, ledger := range result.Ledgers {
		fmt.Fprintf(writer, "%s\t%s\n", ledger.Name, ledger.Description)
	}
	return writer.Flush()
}

type ledgersJSONResponse struct {
	Ledgers []ledgerJSON `json:"ledgers"`
}

type ledgerJSON struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func ledgersJSONResponseFrom(result stewledgers.Result) ledgersJSONResponse {
	ledgers := make([]ledgerJSON, 0, len(result.Ledgers))
	for _, ledger := range result.Ledgers {
		ledgers = append(ledgers, ledgerJSON{
			Name:        ledger.Name,
			Description: ledger.Description,
		})
	}
	return ledgersJSONResponse{Ledgers: ledgers}
}

const ledgersHelp = `List available Stew ledgers.

Stew discovers writable ledgers from ledger specs, excluding the shared model
contract. The default output lists each ledger name and description. Use --json
to print a machine-readable JSON response.

Usage:
  stew ledgers [flags]

Examples:
  stew ledgers
  stew ledgers --json
  stew ledgers --path /path/to/repo

Flags:
  -h, --help          help for ledgers
      --json          Print JSON output
      --path string   Target directory (default ".")
`
