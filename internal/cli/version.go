package cli

import (
	"fmt"

	"github.com/ankitvg/stew/internal/version"
	"github.com/spf13/cobra"
)

func newVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print build version information",
		Long: `Print Stew build metadata.

Use this when debugging which local or released Stew binary is being used. Dev
builds may print placeholder commit/date values unless built with release
linker flags.`,
		Example: `  stew version`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Fprintf(cmd.OutOrStdout(), "version: %s\n", version.Version)
			fmt.Fprintf(cmd.OutOrStdout(), "commit: %s\n", version.Commit)
			fmt.Fprintf(cmd.OutOrStdout(), "date: %s\n", version.Date)
		},
	}
}
