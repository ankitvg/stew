package cli

import (
	"github.com/ankitvg/stew/internal/version"
	"github.com/spf13/cobra"
)

func Execute() error {
	return newRootCmd().Execute()
}

func newRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "stew",
		Short:         "Stew context CLI",
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	cmd.AddCommand(newInitCmd())
	cmd.AddCommand(newFullSpecCmd())
	cmd.AddCommand(newVersionCmd())
	cmd.SetVersionTemplate("{{.Version}}\n")
	cmd.Version = version.Version

	return cmd
}
