package cli

import (
	"fmt"

	"github.com/ankitvg/stew/internal/stewledger"
	"github.com/ankitvg/stew/internal/stewledgercat"
	"github.com/ankitvg/stew/internal/stewledgertail"
)

func runLedger(ctx cliContext, args []string) error {
	if len(args) == 0 || args[0] == "-h" || args[0] == "--help" {
		fmt.Fprint(ctx.out, ledgerHelp)
		return nil
	}

	switch args[0] {
	case "cat":
		return runLedgerCat(ctx, args[1:])
	case "tail":
		return runLedgerTail(ctx, args[1:])
	case "new":
		return runLedgerNew(ctx, args[1:])
	default:
		return fmt.Errorf("unknown ledger command %q", args[0])
	}
}

func runLedgerCat(ctx cliContext, args []string) error {
	if wantsHelp(args) {
		fmt.Fprint(ctx.out, ledgerCatHelp)
		return nil
	}

	var targetPath string

	flags := newFlagSet("ledger cat")
	flags.StringVar(&targetPath, "path", ".", "Target directory")

	positionals, err := parseInterspersedFlags(flags, args, map[string]flagKind{
		"path": stringFlag,
	})
	if err != nil {
		return err
	}
	if err := exactArgs(positionals, 1); err != nil {
		return err
	}

	result, err := stewledgercat.Run(stewledgercat.Options{
		TargetDir: targetPath,
		Ledger:    positionals[0],
	})
	if err != nil {
		return err
	}

	fmt.Fprint(ctx.out, result.Content)
	return nil
}

func runLedgerTail(ctx cliContext, args []string) error {
	if wantsHelp(args) {
		fmt.Fprint(ctx.out, ledgerTailHelp)
		return nil
	}

	var targetPath string
	var limit int

	flags := newFlagSet("ledger tail")
	flags.StringVar(&targetPath, "path", ".", "Target directory")
	flags.IntVar(&limit, "limit", 10, "Number of entries to print")

	positionals, err := parseInterspersedFlags(flags, args, map[string]flagKind{
		"path":  stringFlag,
		"limit": stringFlag,
	})
	if err != nil {
		return err
	}
	if err := exactArgs(positionals, 1); err != nil {
		return err
	}

	result, err := stewledgertail.Run(stewledgertail.Options{
		TargetDir: targetPath,
		Ledger:    positionals[0],
		Limit:     limit,
	})
	if err != nil {
		return err
	}

	fmt.Fprint(ctx.out, result.Content)
	return nil
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

	_, err = stewledger.Run(stewledger.Options{
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

	fmt.Fprintf(ctx.out, "Created ledger %s\n", positionals[0])
	return nil
}

const ledgerHelp = `Manage Stew ledgers.

Stew discovers ledgers from ledger specs. Use "stew ledger new" when you need a
custom append-only ledger beyond the defaults created by init.

Usage:
  stew ledger [command]

Available Commands:
  cat         Print a stew ledger
  new         Create a custom stew ledger
  tail        Print recent ledger entries

Flags:
  -h, --help   help for ledger
`

const ledgerCatHelp = `Print a Stew ledger.

The command writes raw ledger markdown to stdout. Use shell pipes for searching
or filtering, such as "stew ledger cat iterations | grep parser".

Usage:
  stew ledger cat <ledger> [flags]

Examples:
  stew ledger cat iterations
  stew ledger cat decisions --path /path/to/repo
  stew ledger cat iterations | grep "Prompt"

Flags:
  -h, --help          help for cat
      --path string   Target directory (default ".")
`

const ledgerTailHelp = `Print recent Stew ledger entries.

The command writes the last N entries from the ledger to stdout. Output contains
entries only, in file order, with no ledger title or managed marker.

Usage:
  stew ledger tail <ledger> [flags]

Examples:
  stew ledger tail iterations
  stew ledger tail iterations --limit 5
  stew ledger tail decisions --path /path/to/repo --limit 2

Flags:
  -h, --help          help for tail
      --limit int     Number of entries to print (default 10)
      --path string   Target directory (default ".")
`

const ledgerNewHelp = `Create a custom Stew ledger in an initialized repository.

The command creates a ledger and its matching ledger spec. The filesystem
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
