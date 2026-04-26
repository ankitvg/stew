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
