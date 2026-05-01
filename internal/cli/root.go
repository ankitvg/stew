package cli

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/ankitvg/stew/internal/version"
)

func Execute() error {
	return ExecuteWithIO(os.Args[1:], os.Stdin, os.Stdout, os.Stderr)
}

func ExecuteWithIO(args []string, in io.Reader, out io.Writer, errOut io.Writer) error {
	if in == nil {
		in = strings.NewReader("")
	}
	if out == nil {
		out = io.Discard
	}
	if errOut == nil {
		errOut = io.Discard
	}

	ctx := cliContext{
		in:     in,
		out:    out,
		errOut: errOut,
	}
	return runRoot(ctx, args)
}

type cliContext struct {
	in     io.Reader
	out    io.Writer
	errOut io.Writer
}

func runRoot(ctx cliContext, args []string) error {
	if len(args) == 0 {
		fmt.Fprint(ctx.out, rootHelp)
		return nil
	}
	if args[0] == "help" && len(args) == 1 || args[0] == "-h" || args[0] == "--help" {
		fmt.Fprint(ctx.out, rootHelp)
		return nil
	}
	if args[0] == "-v" || args[0] == "--version" {
		fmt.Fprintf(ctx.out, "%s\n", version.Version)
		return nil
	}

	switch args[0] {
	case "help":
		return runHelp(ctx, args[1:])
	case "append":
		return runAppend(ctx, args[1:])
	case "init":
		return runInit(ctx, args[1:])
	case "full-spec":
		return runFullSpec(ctx, args[1:])
	case "ledger":
		return runLedger(ctx, args[1:])
	case "ledgers":
		return runLedgers(ctx, args[1:])
	case "version":
		return runVersion(ctx, args[1:])
	default:
		return fmt.Errorf("unknown command %q", args[0])
	}
}

func runHelp(ctx cliContext, args []string) error {
	if len(args) == 0 {
		fmt.Fprint(ctx.out, rootHelp)
		return nil
	}

	switch args[0] {
	case "append":
		fmt.Fprint(ctx.out, appendHelp)
	case "init":
		fmt.Fprint(ctx.out, initHelp)
	case "full-spec":
		fmt.Fprint(ctx.out, fullSpecHelp)
	case "ledger":
		if len(args) > 1 && args[1] == "new" {
			fmt.Fprint(ctx.out, ledgerNewHelp)
			return nil
		}
		fmt.Fprint(ctx.out, ledgerHelp)
	case "ledgers":
		fmt.Fprint(ctx.out, ledgersHelp)
	case "version":
		fmt.Fprint(ctx.out, versionHelp)
	default:
		return fmt.Errorf("unknown help topic %q", args[0])
	}
	return nil
}

type flagKind int

const (
	stringFlag flagKind = iota
	boolFlag
)

type trackedString struct {
	value   string
	changed bool
}

func (s *trackedString) Set(value string) error {
	s.value = value
	s.changed = true
	return nil
}

func (s *trackedString) String() string {
	return s.value
}

func newFlagSet(name string) *flag.FlagSet {
	flags := flag.NewFlagSet(name, flag.ContinueOnError)
	flags.SetOutput(io.Discard)
	return flags
}

func parseInterspersedFlags(flags *flag.FlagSet, args []string, specs map[string]flagKind) ([]string, error) {
	flagArgs := make([]string, 0, len(args))
	positionals := make([]string, 0, len(args))

	for i := 0; i < len(args); i++ {
		arg := args[i]
		if arg == "--" {
			positionals = append(positionals, args[i+1:]...)
			break
		}
		if !looksLikeFlag(arg) {
			positionals = append(positionals, arg)
			continue
		}

		name := flagName(arg)
		kind, ok := specs[name]
		if !ok {
			return nil, fmt.Errorf("flag provided but not defined: -%s", name)
		}

		flagArgs = append(flagArgs, arg)
		if kind == stringFlag && !strings.Contains(arg, "=") {
			if i+1 >= len(args) {
				return nil, fmt.Errorf("flag needs an argument: -%s", name)
			}
			i++
			flagArgs = append(flagArgs, args[i])
		}
	}

	if err := flags.Parse(flagArgs); err != nil {
		return nil, err
	}
	return positionals, nil
}

func looksLikeFlag(arg string) bool {
	return strings.HasPrefix(arg, "-") && arg != "-"
}

func flagName(arg string) string {
	trimmed := strings.TrimLeft(arg, "-")
	if before, _, found := strings.Cut(trimmed, "="); found {
		return before
	}
	return trimmed
}

func wantsHelp(args []string) bool {
	for _, arg := range args {
		if arg == "-h" || arg == "--help" {
			return true
		}
	}
	return false
}

func exactArgs(args []string, count int) error {
	if len(args) == count {
		return nil
	}
	return fmt.Errorf("accepts %d arg(s), received %d", count, len(args))
}

func requiredFlags(names ...string) error {
	quoted := make([]string, 0, len(names))
	for _, name := range names {
		quoted = append(quoted, fmt.Sprintf("%q", name))
	}
	return fmt.Errorf("required flag(s) %s not set", strings.Join(quoted, ", "))
}

const rootHelp = `Stew maintains append-only markdown ledgers in .stew/.

Use Stew to give agents and humans a durable project memory:
1. Run "stew help" to discover available commands.
2. Run "stew full-spec" to load the repository's ledger contract.
3. Run "stew <command> --help" before using a command you do not know.
4. Use "stew ledger new <name>" to create custom ledgers when needed.
5. Use "stew append <ledger>" to write new entries without editing old ones.

Usage:
  stew [command]

Examples:
  stew init
  stew full-spec
  stew ledger new plans --description "Reasoning artifacts for future work" --threshold "Append when a plan captures durable intent or tradeoffs."
  stew append iterations --help
  printf 'Implemented the change and ran tests.' | stew append iterations --prompt 'Fix bug' --summary 'Fix parser bug'

Available Commands:
  append      Append an entry to a stew ledger
  full-spec   Print the full stew contract from .stew/*.spec.md
  help        Help about any command
  init        Initialize stew context files in a repository
  ledger      Manage stew ledgers
  ledgers     List available ledgers
  version     Print build version information

Flags:
  -h, --help      help for stew
  -v, --version   version for stew

Use "stew [command] --help" for more information about a command.
`
