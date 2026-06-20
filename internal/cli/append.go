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
	var linkFiles trackedStrings
	var jsonOutput bool

	flags := newFlagSet("append")
	flags.BoolVar(&jsonOutput, "json", false, "Print JSON output")
	flags.StringVar(&targetPath, "path", ".", "Target directory")
	flags.Var(&prompt, "prompt", "Originating user prompt")
	flags.Var(&summary, "summary", "Entry summary")
	flags.Var(&message, "message", "Use the given text as the entry body")
	flags.Var(&message, "m", "Use the given text as the entry body")
	flags.Var(&filePath, "file", "Read the entry body from a file")
	flags.Var(&filePath, "F", "Read the entry body from a file")
	flags.Var(&linkFiles, "link-file", "Link the new entry to a repo file")

	positionals, err := parseInterspersedFlags(flags, args, map[string]flagKind{
		"json":      boolFlag,
		"path":      argFlag,
		"prompt":    argFlag,
		"summary":   argFlag,
		"message":   argFlag,
		"m":         argFlag,
		"file":      argFlag,
		"F":         argFlag,
		"link-file": argFlag,
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

	result, err := stewappend.Run(stewappend.Options{
		TargetDir:  targetPath,
		Ledger:     positionals[0],
		Prompt:     prompt.value,
		Summary:    summary.value,
		Message:    message.value,
		MessageSet: message.changed,
		FilePath:   filePath.value,
		LinkFiles:  linkFiles.values,
		Stdin:      stdin,
		StdinIsTTY: stdinIsTTY,
	})
	if err != nil {
		return err
	}

	if jsonOutput {
		response := appendJSONResponse{
			Ledger:   positionals[0],
			EntryRef: result.EntryRef,
		}
		for _, link := range result.Links {
			response.Links = append(response.Links, linkJSON{
				Source:    link.Source,
				Target:    link.Target,
				CreatedAt: link.CreatedAt,
			})
		}
		return encodeJSON(ctx.out, response)
	}
	fmt.Fprintf(ctx.out, "Appended %s\n", positionals[0])
	return nil
}

type appendJSONResponse struct {
	Ledger   string     `json:"ledger"`
	EntryRef string     `json:"entryRef"`
	Links    []linkJSON `json:"links,omitempty"`
}

type trackedStrings struct {
	values []string
}

func (s *trackedStrings) Set(value string) error {
	s.values = append(s.values, value)
	return nil
}

func (s *trackedStrings) String() string {
	return ""
}

const appendHelp = `Append a new entry to a Stew ledger.

The ledger must already be defined by a ledger spec. The reserved ledger name
"stew" is not writable because it names the shared model contract.

Stew owns the entry header, timestamp, prompt field, and separator. Supply the
ledger name, the originating prompt, a short summary, and exactly one body
source: piped stdin, -m/--message, or -F/--file.

Use "stew full-spec" to read the repository's ledger rules before appending.
Use "stew append <ledger> --help" when you need the command contract.
Use --json when an agent or script needs the created entry ref.
Use --link-file to associate the new entry with one or more repo files.

Usage:
  stew append <ledger> [flags]

Examples:
  printf 'Implemented the change and ran go test ./...' | stew append iterations --prompt 'Add append command' --summary 'Implement append command'
  stew append iterations --prompt 'Small fix' --summary 'Record small fix' -m 'Updated validation and ran tests.'
  stew append decisions --prompt 'Choose storage model' --summary 'Use append-only ledgers' -F decision-entry.md
  stew append iterations --prompt 'Fix parser' --summary 'Fix parser' -m 'Updated parser.' --link-file internal/parser.go

Flags:
  -F, --file string      Read the entry body from a file
  -h, --help             help for append
      --json             Print JSON output
      --link-file value  Link the new entry to a repo file
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
