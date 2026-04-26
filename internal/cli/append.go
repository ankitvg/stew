package cli

import (
	"fmt"
	"io"
	"os"

	"github.com/ankitvg/stew/internal/stewappend"
	"github.com/spf13/cobra"
)

func newAppendCmd() *cobra.Command {
	var targetPath string
	var prompt string
	var summary string
	var message string
	var filePath string

	cmd := &cobra.Command{
		Use:   "append <ledger>",
		Short: "Append an entry to a stew ledger",
		Long: `Append a new entry to .stew/<ledger>.md.

The ledger must already be defined by .stew/<ledger>.spec.md. The reserved
ledger name "stew" is not writable because .stew/stew.spec.md defines the
shared model contract.

Stew owns the entry header, timestamp, prompt field, and separator. Supply the
ledger name, the originating prompt, a short summary, and exactly one body
source: piped stdin, -m/--message, or -F/--file.

Use "stew full-spec" to read the repository's ledger rules before appending.
Use "stew append <ledger> --help" when you need the command contract.`,
		Example: `  printf 'Implemented the change and ran go test ./...' | stew append iterations --prompt 'Add append command' --summary 'Implement append command'
  stew append iterations --prompt 'Small fix' --summary 'Record small fix' -m 'Updated validation and ran tests.'
  stew append decisions --prompt 'Choose storage model' --summary 'Use append-only ledgers' -F decision-entry.md`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			messageSet := cmd.Flags().Changed("message")
			fileSet := cmd.Flags().Changed("file")
			var stdin io.Reader
			var stdinIsTTY func() bool
			if !messageSet && !fileSet {
				stdin = cmd.InOrStdin()
				stdinIsTTY = func() bool { return readerIsTerminal(cmd.InOrStdin()) }
			}

			result, err := stewappend.Run(stewappend.Options{
				TargetDir:  targetPath,
				Ledger:     args[0],
				Prompt:     prompt,
				Summary:    summary,
				Message:    message,
				MessageSet: messageSet,
				FilePath:   filePath,
				Stdin:      stdin,
				StdinIsTTY: stdinIsTTY,
			})
			if err != nil {
				return err
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Appended %s\n", result.LedgerPath)
			return nil
		},
	}

	cmd.Flags().StringVar(&targetPath, "path", ".", "Target directory")
	cmd.Flags().StringVar(&prompt, "prompt", "", "Originating user prompt")
	cmd.Flags().StringVar(&summary, "summary", "", "Entry summary")
	cmd.Flags().StringVarP(&message, "message", "m", "", "Use the given text as the entry body")
	cmd.Flags().StringVarP(&filePath, "file", "F", "", "Read the entry body from a file")
	_ = cmd.MarkFlagRequired("prompt")
	_ = cmd.MarkFlagRequired("summary")

	return cmd
}

func readerIsTerminal(reader io.Reader) bool {
	file, ok := reader.(*os.File)
	if !ok {
		return false
	}
	info, err := file.Stat()
	if err != nil {
		return false
	}
	return info.Mode()&os.ModeCharDevice != 0
}
