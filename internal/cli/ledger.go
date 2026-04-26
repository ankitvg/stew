package cli

import (
	"fmt"

	"github.com/ankitvg/stew/internal/stewledger"
	"github.com/spf13/cobra"
)

func newLedgerCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ledger",
		Short: "Manage stew ledgers",
		Long: `Manage Stew ledgers.

Stew discovers ledgers from .stew/*.spec.md files. Use "stew ledger new" when
you need a custom append-only ledger beyond the defaults created by init.`,
	}

	cmd.AddCommand(newLedgerNewCmd())
	return cmd
}

func newLedgerNewCmd() *cobra.Command {
	var targetPath string
	var description string
	var threshold string
	var quiet bool

	cmd := &cobra.Command{
		Use:   "new <name>",
		Short: "Create a custom stew ledger",
		Long: `Create a custom Stew ledger in an initialized repository.

The command writes .stew/<name>.md and .stew/<name>.spec.md. The filesystem
remains the source of truth: no registry, config, or generated code is updated.

Names must use lowercase ASCII letters, digits, and single hyphens only. Supply
--description and --threshold to make the generated spec useful immediately; if
omitted, the spec contains TODO guidance that is visible in stew full-spec.`,
		Example: `  stew ledger new plans --description "Reasoning artifacts for future work" --threshold "Append when a plan captures durable intent or tradeoffs."
  stew ledger new experiments --path /path/to/repo
  stew ledger new reviews --quiet`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			result, err := stewledger.Run(stewledger.Options{
				TargetDir:   targetPath,
				Name:        args[0],
				Description: description,
				Threshold:   threshold,
			})
			if err != nil {
				return err
			}
			if quiet {
				return nil
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Created %s\n", result.LedgerPath)
			fmt.Fprintf(cmd.OutOrStdout(), "Created %s\n", result.SpecPath)
			return nil
		},
	}

	cmd.Flags().StringVar(&targetPath, "path", ".", "Target directory")
	cmd.Flags().StringVar(&description, "description", "", "Ledger description for the generated spec")
	cmd.Flags().StringVar(&threshold, "threshold", "", "Guidance for when to append entries")
	cmd.Flags().BoolVar(&quiet, "quiet", false, "Suppress standard output")

	return cmd
}
