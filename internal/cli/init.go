package cli

import (
	"fmt"
	"sort"

	"github.com/ankitvg/stew/internal/stewinit"
	"github.com/spf13/cobra"
)

func newInitCmd() *cobra.Command {
	var targetPath string
	var noAgentsMD bool
	var quiet bool

	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize stew context files in a repository",
		Long: `Initialize Stew in a repository.

This command creates the .stew/ directory, the default ledger/spec files, and
the managed Stew block in AGENTS.md. It is safe to run more than once: existing
ledger files are preserved, and the managed AGENTS.md block is inserted or
updated without replacing content outside the block.

After init, agents should use "stew help" for CLI discovery. Run "stew full-spec"
to load the repository's ledger contract.`,
		Example: `  stew init
  stew init --path /path/to/repo
  stew init --no-agents-md
  stew init --quiet`,
		RunE: func(cmd *cobra.Command, args []string) error {
			result, err := stewinit.Run(stewinit.Options{
				TargetDir:  targetPath,
				NoAgentsMD: noAgentsMD,
			})
			if err != nil {
				return err
			}

			printWarnings(cmd, result.Warnings)
			if quiet {
				return nil
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Initialized stew in %s\n", result.TargetDir)

			paths := make([]string, 0, len(result.FileStatuses))
			for path := range result.FileStatuses {
				paths = append(paths, path)
			}
			sort.Strings(paths)
			for _, path := range paths {
				fmt.Fprintf(cmd.OutOrStdout(), "[%s] %s\n", result.FileStatuses[path], path)
			}

			if !noAgentsMD {
				fmt.Fprintf(cmd.OutOrStdout(), "[%s] AGENTS.md\n", result.AgentsStatus)
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&targetPath, "path", ".", "Target directory")
	cmd.Flags().BoolVar(&noAgentsMD, "no-agents-md", false, "Skip AGENTS.md managed block updates")
	cmd.Flags().BoolVar(&quiet, "quiet", false, "Suppress standard output")

	return cmd
}

func printWarnings(cmd *cobra.Command, warnings []string) {
	for _, warning := range warnings {
		fmt.Fprintf(cmd.ErrOrStderr(), "warning: %s\n", warning)
	}
}
