package cli

import (
	"fmt"

	"github.com/ankitvg/stew/internal/stewledger"
)

func runLedger(ctx cliContext, args []string) error {
	if len(args) == 0 || args[0] == "-h" || args[0] == "--help" {
		fmt.Fprint(ctx.out, ledgerHelp)
		return nil
	}

	if args[0] != "new" {
		return fmt.Errorf("unknown ledger command %q", args[0])
	}
	return runLedgerNew(ctx, args[1:])
}

func runLedgerNew(ctx cliContext, args []string) error {
	if wantsHelp(args) {
		fmt.Fprint(ctx.out, ledgerNewHelp)
		return nil
	}

	var description string
	var threshold string
	var quiet bool

	flags := newFlagSet("ledger new")
	flags.StringVar(&description, "description", "", "Ledger description for the generated spec")
	flags.StringVar(&threshold, "threshold", "", "Guidance for when to append entries")
	flags.BoolVar(&quiet, "quiet", false, "Suppress standard output")

	positionals, err := parseInterspersedFlags(flags, args, map[string]flagKind{
		"description": stringFlag,
		"threshold":   stringFlag,
		"quiet":       boolFlag,
	})
	if err != nil {
		return err
	}
	if err := exactArgs(positionals, 1); err != nil {
		return err
	}

	result, err := stewledger.Run(stewledger.Options{
		Name:        positionals[0],
		Description: description,
		Threshold:   threshold,
	})
	if err != nil {
		return err
	}
	if quiet {
		return nil
	}

	fmt.Fprintf(ctx.out, "Created %s\n", result.LedgerPath)
	fmt.Fprintf(ctx.out, "Created %s\n", result.SpecPath)
	return nil
}

const ledgerHelp = `Manage Stew ledgers.

Stew discovers ledgers from .stew/*.spec.md files. Use "stew ledger new" when
you need a custom append-only ledger beyond the defaults created by init.

Usage:
  stew ledger [command]

Available Commands:
  new         Create a custom stew ledger

Flags:
  -h, --help   help for ledger
`

const ledgerNewHelp = `Create a custom Stew ledger in an initialized repository.

The command writes .stew/<name>.md and .stew/<name>.spec.md. The filesystem
remains the source of truth: no registry, config, or generated code is updated.

Names must use lowercase ASCII letters, digits, and single hyphens only. Supply
--description and --threshold to make the generated spec useful immediately; if
omitted, the spec contains TODO guidance that is visible in stew full-spec.

Usage:
  stew ledger new <name> [flags]

Examples:
  stew ledger new plans --description "Reasoning artifacts for future work" --threshold "Append when a plan captures durable intent or tradeoffs."
  stew ledger new experiments
  stew ledger new reviews --quiet

Flags:
      --description string   Ledger description for the generated spec
  -h, --help                 help for new
      --quiet                Suppress standard output
      --threshold string     Guidance for when to append entries
`
