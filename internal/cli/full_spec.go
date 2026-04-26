package cli

import (
	"fmt"

	"github.com/ankitvg/stew/internal/fullspec"
	"github.com/spf13/cobra"
)

func newFullSpecCmd() *cobra.Command {
	var targetPath string

	cmd := &cobra.Command{
		Use:   "full-spec",
		Short: "Print the full stew contract from .stew/*.spec.md",
		Long: `Print the full Stew contract for the target repository.

The output starts with .stew/stew.spec.md, then includes every other
.stew/*.spec.md file in deterministic order. Agents should run this before
writing ledger entries so they know the repository's append-only rules and each
ledger's expected entry shape.`,
		Example: `  stew full-spec
  stew full-spec --path /path/to/repo
  stew full-spec | sed -n '1,120p'`,
		RunE: func(cmd *cobra.Command, args []string) error {
			result, err := fullspec.Load(targetPath)
			if err != nil {
				return err
			}
			fmt.Fprint(cmd.OutOrStdout(), result.Content)
			return nil
		},
	}

	cmd.Flags().StringVar(&targetPath, "path", ".", "Target directory")

	return cmd
}
