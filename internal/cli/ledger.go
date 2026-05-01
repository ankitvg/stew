package cli

import (
	"fmt"
	"strings"

	"github.com/ankitvg/stew/internal/stewledger"
	"github.com/ankitvg/stew/internal/stewledgerall"
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

	var all bool
	var targetPath string

	flags := newFlagSet("ledger cat")
	flags.BoolVar(&all, "all", false, "Print all ledgers")
	flags.StringVar(&targetPath, "path", ".", "Target directory")

	positionals, err := parseInterspersedFlags(flags, args, map[string]flagKind{
		"all":  boolFlag,
		"path": argFlag,
	})
	if err != nil {
		return err
	}
	if all {
		if len(positionals) > 0 {
			return fmt.Errorf("cannot use --all with a ledger argument")
		}
		result, err := stewledgerall.Cat(stewledgerall.Options{
			TargetDir: targetPath,
		})
		if err != nil {
			return err
		}
		fmt.Fprint(ctx.out, renderLedgerSections(result.Sections))
		return nil
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

	var all bool
	var targetPath string
	var limit int

	flags := newFlagSet("ledger tail")
	flags.BoolVar(&all, "all", false, "Print entries from all ledgers")
	flags.StringVar(&targetPath, "path", ".", "Target directory")
	flags.IntVar(&limit, "limit", 10, "Number of entries to print")

	positionals, err := parseInterspersedFlags(flags, args, map[string]flagKind{
		"all":   boolFlag,
		"path":  argFlag,
		"limit": argFlag,
	})
	if err != nil {
		return err
	}
	if all {
		if len(positionals) > 0 {
			return fmt.Errorf("cannot use --all with a ledger argument")
		}
		result, err := stewledgerall.Tail(stewledgerall.Options{
			TargetDir: targetPath,
			Limit:     limit,
		})
		if err != nil {
			return err
		}
		fmt.Fprint(ctx.out, renderLedgerSections(result.Sections))
		return nil
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
		"description": argFlag,
		"threshold":   argFlag,
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

func renderLedgerSections(sections []stewledgerall.Section) string {
	var builder strings.Builder
	for i, section := range sections {
		if i > 0 {
			builder.WriteString("\n")
		}
		builder.WriteString("# ")
		builder.WriteString(section.Name)
		builder.WriteString("\n\n")
		builder.WriteString(section.Content)
		if section.Content != "" && !strings.HasSuffix(section.Content, "\n") {
			builder.WriteString("\n")
		}
	}
	return builder.String()
}

const ledgerHelp = `Manage Stew ledgers.

Stew discovers ledgers from ledger specs. Use "stew ledger new" when you need a
custom append-only ledger beyond the defaults created by init.

Usage:
  stew ledger [command]

Available Commands:
  cat         Print one or all stew ledgers
  new         Create a custom stew ledger
  tail        Print recent ledger entries

Flags:
  -h, --help   help for ledger
`

const ledgerCatHelp = `Print Stew ledger content.

For one ledger, the command writes raw ledger markdown to stdout. With --all,
the command writes each ledger under a ledger-name section. Use shell pipes for
searching or filtering, such as "stew ledger cat iterations | grep parser".

Usage:
  stew ledger cat <ledger> [flags]
  stew ledger cat --all [flags]

Examples:
  stew ledger cat iterations
  stew ledger cat --all
  stew ledger cat decisions --path /path/to/repo
  stew ledger cat iterations | grep "Prompt"

Flags:
      --all           Print all ledgers
  -h, --help          help for cat
      --path string   Target directory (default ".")
`

const ledgerTailHelp = `Print recent Stew ledger entries.

For one ledger, the command writes the last N entries to stdout. With --all,
the command writes each ledger's last N entries under a ledger-name section.
Entries remain in file order.

Usage:
  stew ledger tail <ledger> [flags]
  stew ledger tail --all [flags]

Examples:
  stew ledger tail iterations
  stew ledger tail iterations --limit 5
  stew ledger tail --all --limit 5
  stew ledger tail decisions --path /path/to/repo --limit 2

Flags:
      --all           Print entries from all ledgers
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
