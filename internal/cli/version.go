package cli

import (
	"fmt"

	"github.com/ankitvg/stew/internal/version"
)

func runVersion(ctx cliContext, args []string) error {
	if wantsHelp(args) {
		fmt.Fprint(ctx.out, versionHelp)
		return nil
	}
	if err := exactArgs(args, 0); err != nil {
		return err
	}

	fmt.Fprintf(ctx.out, "version: %s\n", version.Version)
	fmt.Fprintf(ctx.out, "commit: %s\n", version.Commit)
	fmt.Fprintf(ctx.out, "date: %s\n", version.Date)
	return nil
}

const versionHelp = `Print Stew build metadata.

Use this when debugging which local or released Stew binary is being used. Dev
builds may print placeholder commit/date values unless built with release
linker flags.

Usage:
  stew version

Examples:
  stew version
`
