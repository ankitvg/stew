package cli

import (
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

	flags := newFlagSet("ledgers")
	flags.StringVar(&targetPath, "path", ".", "Target directory")

	positionals, err := parseInterspersedFlags(flags, args, map[string]flagKind{
		"path": stringFlag,
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

	writer := tabwriter.NewWriter(ctx.out, 0, 0, 2, ' ', 0)
	for _, ledger := range result.Ledgers {
		fmt.Fprintf(writer, "%s\t%s\t%s\n", ledger.Name, ledger.LedgerPath, ledger.Description)
	}
	return writer.Flush()
}

const ledgersHelp = `List available Stew ledgers.

Stew discovers writable ledgers from .stew/*.spec.md files, excluding the
shared .stew/stew.spec.md model contract. The output lists each ledger name,
entry file path, and description.

Usage:
  stew ledgers [flags]

Examples:
  stew ledgers
  stew ledgers --path /path/to/repo

Flags:
  -h, --help          help for ledgers
      --path string   Target directory (default ".")
`
