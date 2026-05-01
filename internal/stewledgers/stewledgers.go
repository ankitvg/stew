package stewledgers

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

var (
	ErrMissingStewDir = errors.New("stew is not initialized")
	ErrNoLedgers      = errors.New("no writable ledgers found")
)

type Options struct {
	TargetDir string
}

type Ledger struct {
	Name        string
	LedgerPath  string
	Description string
}

type Result struct {
	TargetDir string
	Ledgers   []Ledger
}

func List(opts Options) (Result, error) {
	targetDir := opts.TargetDir
	if strings.TrimSpace(targetDir) == "" {
		targetDir = "."
	}

	absTarget, err := filepath.Abs(targetDir)
	if err != nil {
		return Result{}, fmt.Errorf("resolve target path: %w", err)
	}

	stewDir := filepath.Join(absTarget, ".stew")
	entries, err := os.ReadDir(stewDir)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return Result{}, fmt.Errorf("%w in %s", ErrMissingStewDir, absTarget)
		}
		return Result{}, fmt.Errorf("read stew metadata: %w", err)
	}

	ledgers := make([]Ledger, 0)
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		fileName := entry.Name()
		if fileName == "stew.spec.md" || !strings.HasSuffix(fileName, ".spec.md") {
			continue
		}

		name := strings.TrimSuffix(fileName, ".spec.md")
		specRel := filepath.Join(".stew", fileName)
		specPath := filepath.Join(absTarget, specRel)
		bytes, err := os.ReadFile(specPath)
		if err != nil {
			return Result{}, fmt.Errorf("read ledger spec for %q: %w", name, err)
		}

		ledgers = append(ledgers, Ledger{
			Name:        name,
			LedgerPath:  filepath.Join(".stew", name+".md"),
			Description: extractDescription(string(bytes)),
		})
	}

	if len(ledgers) == 0 {
		return Result{}, fmt.Errorf("%w in %s", ErrNoLedgers, absTarget)
	}

	sort.Slice(ledgers, func(i, j int) bool {
		return ledgers[i].Name < ledgers[j].Name
	})

	return Result{
		TargetDir: absTarget,
		Ledgers:   ledgers,
	}, nil
}

func extractDescription(content string) string {
	lines := strings.Split(strings.ReplaceAll(content, "\r\n", "\n"), "\n")

	start := 0
	for i, raw := range lines {
		if strings.HasPrefix(strings.TrimSpace(raw), "# ") {
			start = i + 1
			break
		}
	}

	paragraph := make([]string, 0)
	for _, raw := range lines[start:] {
		line := strings.TrimSpace(raw)
		if line == "" {
			if len(paragraph) > 0 {
				break
			}
			continue
		}
		if startsNonProseBlock(line) {
			if len(paragraph) > 0 {
				break
			}
			continue
		}
		paragraph = append(paragraph, line)
	}

	return strings.Join(paragraph, " ")
}

func startsNonProseBlock(line string) bool {
	return strings.HasPrefix(line, "#") ||
		strings.HasPrefix(line, "<!--") ||
		strings.HasPrefix(line, "```") ||
		strings.HasPrefix(line, "- ") ||
		strings.HasPrefix(line, "* ") ||
		strings.HasPrefix(line, ">")
}
