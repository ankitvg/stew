package cli

import (
	"fmt"

	"github.com/ankitvg/stew/internal/stewlink"
	"github.com/ankitvg/stew/internal/stewref"
)

func runLink(ctx cliContext, args []string) error {
	if len(args) == 0 || args[0] == "-h" || args[0] == "--help" {
		fmt.Fprint(ctx.out, linkHelp)
		return nil
	}

	switch args[0] {
	case "list":
		return runLinkList(ctx, args[1:])
	default:
		return fmt.Errorf("unknown link command %q", args[0])
	}
}

func runLinkList(ctx cliContext, args []string) error {
	if wantsHelp(args) {
		fmt.Fprint(ctx.out, linkListHelp)
		return nil
	}

	var jsonOutput bool
	var targetPath string

	flags := newFlagSet("link list")
	flags.BoolVar(&jsonOutput, "json", false, "Print JSON output")
	flags.StringVar(&targetPath, "path", ".", "Target directory")

	positionals, err := parseInterspersedFlags(flags, args, map[string]flagKind{
		"json": boolFlag,
		"path": argFlag,
	})
	if err != nil {
		return err
	}
	if err := exactArgs(positionals, 1); err != nil {
		return err
	}

	ref, err := stewref.Parse(positionals[0])
	if err != nil {
		return err
	}
	result, err := stewlink.List(stewlink.ListOptions{
		TargetDir: targetPath,
		Ref:       ref,
	})
	if err != nil {
		return err
	}

	if jsonOutput {
		return encodeJSON(ctx.out, linkListJSONResponseFrom(result))
	}
	for _, link := range result.Links {
		fmt.Fprintf(ctx.out, "%s -> %s\n", link.Source, link.Target)
	}
	return nil
}

type linkListJSONResponse struct {
	Ref   string     `json:"ref"`
	Links []linkJSON `json:"links"`
}

type linkJSON struct {
	Source    string `json:"source"`
	Target    string `json:"target"`
	CreatedAt string `json:"createdAt"`
}

func linkListJSONResponseFrom(result stewlink.ListResult) linkListJSONResponse {
	links := make([]linkJSON, 0, len(result.Links))
	for _, link := range result.Links {
		links = append(links, linkJSON{
			Source:    link.Source,
			Target:    link.Target,
			CreatedAt: link.CreatedAt,
		})
	}
	return linkListJSONResponse{
		Ref:   result.Ref,
		Links: links,
	}
}

const linkHelp = `Manage Stew links.

Links are append-only relationships between refs. V1 links have a source ref
and a target ref, without a kind. Use append --link-file to create links.

Usage:
  stew link [command]

Available Commands:
  list        List links for a ref

Flags:
  -h, --help   help for link
`

const linkListHelp = `List Stew links for a ref.

The command prints links where the given ref is either the source or target.
Default output prints one "source -> target" line per link. Use --json for
machine-readable output.

Usage:
  stew link list <ref> [flags]

Examples:
  stew link list file:internal/cli/ledger.go
  stew link list entry:iterations/2026-06-20T211727Z-vwli67-plumb-refs-into-surfaces.md --json

Flags:
  -h, --help          help for list
      --json          Print JSON output
      --path string   Target directory (default ".")
`
