package cli

import (
	"fmt"
	"io"
	"sort"

	"github.com/ankitvg/stew/internal/stewinit"
)

func runInit(ctx cliContext, args []string) error {
	if wantsHelp(args) {
		fmt.Fprint(ctx.out, initHelp)
		return nil
	}

	var targetPath string
	var noAgentsMD bool
	var quiet bool

	flags := newFlagSet("init")
	flags.StringVar(&targetPath, "path", ".", "Target directory")
	flags.BoolVar(&noAgentsMD, "no-agents-md", false, "Skip AGENTS.md managed block updates")
	flags.BoolVar(&quiet, "quiet", false, "Suppress standard output")

	positionals, err := parseInterspersedFlags(flags, args, map[string]flagKind{
		"path":         stringFlag,
		"no-agents-md": boolFlag,
		"quiet":        boolFlag,
	})
	if err != nil {
		return err
	}
	if err := exactArgs(positionals, 0); err != nil {
		return err
	}

	result, err := stewinit.Run(stewinit.Options{
		TargetDir:  targetPath,
		NoAgentsMD: noAgentsMD,
	})
	if err != nil {
		return err
	}

	printWarnings(ctx.errOut, result.Warnings)
	if quiet {
		return nil
	}

	fmt.Fprintf(ctx.out, "Initialized stew in %s\n", result.TargetDir)

	paths := make([]string, 0, len(result.FileStatuses))
	for path := range result.FileStatuses {
		paths = append(paths, path)
	}
	sort.Strings(paths)
	for _, path := range paths {
		fmt.Fprintf(ctx.out, "[%s] %s\n", result.FileStatuses[path], path)
	}

	if !noAgentsMD {
		fmt.Fprintf(ctx.out, "[%s] AGENTS.md\n", result.AgentsStatus)
	}

	return nil
}

const initHelp = `Initialize Stew in a repository.

This command creates the .stew/ directory, the default ledger/spec files, and
the managed Stew block in AGENTS.md. It is safe to run more than once: existing
ledger files are preserved, and the managed AGENTS.md block is inserted or
updated without replacing content outside the block.

After init, agents should use "stew help" for CLI discovery. Run "stew full-spec"
to load the repository's ledger contract.

Usage:
  stew init [flags]

Examples:
  stew init
  stew init --path /path/to/repo
  stew init --no-agents-md
  stew init --quiet

Flags:
  -h, --help            help for init
      --no-agents-md    Skip AGENTS.md managed block updates
      --path string     Target directory (default ".")
      --quiet           Suppress standard output
`

func printWarnings(errOut io.Writer, warnings []string) {
	for _, warning := range warnings {
		fmt.Fprintf(errOut, "warning: %s\n", warning)
	}
}
