package cli

import (
	"fmt"

	"github.com/ankitvg/stew/internal/fullspec"
)

func runFullSpec(ctx cliContext, args []string) error {
	if wantsHelp(args) {
		fmt.Fprint(ctx.out, fullSpecHelp)
		return nil
	}

	var targetPath string

	flags := newFlagSet("full-spec")
	flags.StringVar(&targetPath, "path", ".", "Target directory")

	positionals, err := parseInterspersedFlags(flags, args, map[string]flagKind{
		"path": stringFlag,
	})
	if err != nil {
		return err
	}
	if err := exactArgs(positionals, 0); err != nil {
		return err
	}

	result, err := fullspec.Load(targetPath)
	if err != nil {
		return err
	}
	fmt.Fprint(ctx.out, result.Content)
	return nil
}

const fullSpecHelp = `Print the full Stew contract for the target repository.

The output starts with .stew/stew.spec.md, then includes every other
.stew/*.spec.md file in deterministic order. Agents should run this before
writing ledger entries so they know the repository's append-only rules and each
ledger's expected entry shape.

Usage:
  stew full-spec [flags]

Examples:
  stew full-spec
  stew full-spec --path /path/to/repo
  stew full-spec | sed -n '1,120p'

Flags:
  -h, --help          help for full-spec
      --path string   Target directory (default ".")
`
