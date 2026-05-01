package cli

import (
	"fmt"
	"io"
	"os"

	"github.com/ankitvg/stew/internal/stewappend"
)

func runAppend(ctx cliContext, args []string) error {
	if wantsHelp(args) {
		fmt.Fprint(ctx.out, appendHelp)
		return nil
	}

	var targetPath string
	var prompt trackedString
	var summary trackedString
	var message trackedString
	var filePath trackedString

	flags := newFlagSet("append")
	flags.StringVar(&targetPath, "path", ".", "Target directory")
	flags.Var(&prompt, "prompt", "Originating user prompt")
	flags.Var(&summary, "summary", "Entry summary")
	flags.Var(&message, "message", "Use the given text as the entry body")
	flags.Var(&message, "m", "Use the given text as the entry body")
	flags.Var(&filePath, "file", "Read the entry body from a file")
	flags.Var(&filePath, "F", "Read the entry body from a file")

	positionals, err := parseInterspersedFlags(flags, args, map[string]flagKind{
		"path":    argFlag,
		"prompt":  argFlag,
		"summary": argFlag,
		"message": argFlag,
		"m":       argFlag,
		"file":    argFlag,
		"F":       argFlag,
	})
	if err != nil {
		return err
	}
	if err := exactArgs(positionals, 1); err != nil {
		return err
	}

	var missing []string
	if !prompt.changed {
		missing = append(missing, "prompt")
	}
	if !summary.changed {
		missing = append(missing, "summary")
	}
	if len(missing) > 0 {
		return requiredFlags(missing...)
	}

	var stdin io.Reader
	var stdinIsTTY func() bool
	if !message.changed && !filePath.changed {
		stdin = ctx.in
		stdinIsTTY = func() bool { return readerIsTerminal(ctx.in) }
	}

	_, err = stewappend.Run(stewappend.Options{
		TargetDir:  targetPath,
		Ledger:     positionals[0],
		Prompt:     prompt.value,
		Summary:    summary.value,
		Message:    message.value,
		MessageSet: message.changed,
		FilePath:   filePath.value,
		Stdin:      stdin,
		StdinIsTTY: stdinIsTTY,
	})
	if err != nil {
		return err
	}

	fmt.Fprintf(ctx.out, "Appended %s\n", positionals[0])
	return nil
}

const appendHelp = `Append a new entry to a Stew ledger.

The ledger must already be defined by a ledger spec. The reserved ledger name
"stew" is not writable because it names the shared model contract.

Stew owns the entry header, timestamp, prompt field, and separator. Supply the
ledger name, the originating prompt, a short summary, and exactly one body
source: piped stdin, -m/--message, or -F/--file.

Use "stew full-spec" to read the repository's ledger rules before appending.
Use "stew append <ledger> --help" when you need the command contract.

Usage:
  stew append <ledger> [flags]

Examples:
  printf 'Implemented the change and ran go test ./...' | stew append iterations --prompt 'Add append command' --summary 'Implement append command'
  stew append iterations --prompt 'Small fix' --summary 'Record small fix' -m 'Updated validation and ran tests.'
  stew append decisions --prompt 'Choose storage model' --summary 'Use append-only ledgers' -F decision-entry.md

Flags:
  -F, --file string      Read the entry body from a file
  -h, --help             help for append
  -m, --message string   Use the given text as the entry body
      --path string      Target directory (default ".")
      --prompt string    Originating user prompt
      --summary string   Entry summary
`

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
