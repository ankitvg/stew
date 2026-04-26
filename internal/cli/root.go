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
		Use:   "stew",
		Short: "Stew context CLI",
		Long: `Stew maintains append-only markdown ledgers in .stew/.

Use Stew to give agents and humans a durable project memory:
1. Run "stew help" to discover available commands.
2. Run "stew full-spec" to load the repository's ledger contract.
3. Run "stew <command> --help" before using a command you do not know.
4. Use "stew ledger new <name>" to create custom ledgers when needed.
5. Use "stew append <ledger>" to write new entries without editing old ones.`,
		Example: `  stew init
  stew full-spec
  stew ledger new plans --description "Reasoning artifacts for future work" --threshold "Append when a plan captures durable intent or tradeoffs."
  stew append iterations --help
  printf 'Implemented the change and ran tests.' | stew append iterations --prompt 'Fix bug' --summary 'Fix parser bug'`,
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	cmd.AddCommand(newAppendCmd())
	cmd.AddCommand(newInitCmd())
	cmd.AddCommand(newFullSpecCmd())
	cmd.AddCommand(newLedgerCmd())
	cmd.AddCommand(newVersionCmd())
	cmd.SetVersionTemplate("{{.Version}}\n")
	cmd.Version = version.Version

	return cmd
}
