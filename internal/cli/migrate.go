package cli

import (
	"fmt"

	"github.com/ankitvg/stew/internal/stewmigrateatomic"
)

func runMigrate(ctx cliContext, args []string) error {
	if len(args) == 0 || args[0] == "-h" || args[0] == "--help" {
		fmt.Fprint(ctx.out, migrateHelp)
		return nil
	}

	switch args[0] {
	case "atomic-entries":
		return runMigrateAtomicEntries(ctx, args[1:])
	default:
		return fmt.Errorf("unknown migrate command %q", args[0])
	}
}

func runMigrateAtomicEntries(ctx cliContext, args []string) error {
	if wantsHelp(args) {
		fmt.Fprint(ctx.out, migrateAtomicEntriesHelp)
		return nil
	}

	var targetPath string
	var dryRun bool

	flags := newFlagSet("migrate atomic-entries")
	flags.StringVar(&targetPath, "path", ".", "Target directory")
	flags.BoolVar(&dryRun, "dry-run", false, "Preview migration without writing files")

	positionals, err := parseInterspersedFlags(flags, args, map[string]flagKind{
		"path":    argFlag,
		"dry-run": boolFlag,
	})
	if err != nil {
		return err
	}
	if err := exactArgs(positionals, 0); err != nil {
		return err
	}

	result, err := stewmigrateatomic.Run(stewmigrateatomic.Options{
		TargetDir: targetPath,
		DryRun:    dryRun,
	})
	if err != nil {
		return err
	}

	if dryRun {
		fmt.Fprintf(ctx.out, "Would migrate %d ledger(s), write %d entries, and remove %d legacy ledger file(s)\n", result.WouldMigrate, result.WouldWrite, result.WouldRemove)
		return nil
	}
	fmt.Fprintf(ctx.out, "Migrated %d ledger(s), wrote %d entries, and removed %d legacy ledger file(s)\n", result.LedgerCount, result.EntryCount, result.RemovedCount)
	return nil
}

const migrateHelp = `Migrate Stew metadata between storage formats.

Usage:
  stew migrate [command]

Available Commands:
  atomic-entries   Split monolithic ledger files into atomic entry files

Flags:
  -h, --help   help for migrate
`

const migrateAtomicEntriesHelp = `Split monolithic Stew ledgers into atomic entry files.

The command reads existing .stew/<ledger>.md files, writes entries under
.stew/ledgers/<ledger>/, verifies the written entries, then removes the legacy
ledger files. Use --dry-run to preview counts without writing or removing files.

Usage:
  stew migrate atomic-entries [flags]

Examples:
  stew migrate atomic-entries
  stew migrate atomic-entries --dry-run
  stew migrate atomic-entries --path /path/to/repo

Flags:
      --dry-run       Preview migration without writing files
  -h, --help          help for atomic-entries
      --path string   Target directory (default ".")
`
