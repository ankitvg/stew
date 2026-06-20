package stewledgercat

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

var (
	ErrInvalidLedger = errors.New("invalid ledger")
	ErrUnknownLedger = errors.New("unknown ledger")
	ErrMissingLedger = errors.New("missing ledger")
)

type Options struct {
	TargetDir string
	Ledger    string
}

type Result struct {
	TargetDir  string
	LedgerPath string
	Content    string
	EntryFiles []EntryFile
}

type EntryFile struct {
	Name    string
	RelPath string
	Content string
}

func Run(opts Options) (Result, error) {
	if opts.TargetDir == "" {
		opts.TargetDir = "."
	}

	targetDir, err := filepath.Abs(opts.TargetDir)
	if err != nil {
		return Result{}, fmt.Errorf("resolve target path: %w", err)
	}
	info, err := os.Stat(targetDir)
	if err != nil {
		return Result{}, fmt.Errorf("stat target path: %w", err)
	}
	if !info.IsDir() {
		return Result{}, fmt.Errorf("target path is not a directory: %s", targetDir)
	}

	ledger, err := validateLedgerName(opts.Ledger)
	if err != nil {
		return Result{}, err
	}

	stewDir := filepath.Join(targetDir, ".stew")
	specRel := filepath.Join(".stew", ledger+".spec.md")
	specPath := filepath.Join(targetDir, specRel)
	specInfo, err := os.Stat(specPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return Result{}, fmt.Errorf("%w: ledger %q is not defined", ErrUnknownLedger, ledger)
		}
		return Result{}, fmt.Errorf("stat ledger spec: %w", err)
	}
	if specInfo.IsDir() {
		return Result{}, fmt.Errorf("%w: ledger %q is not defined", ErrUnknownLedger, ledger)
	}

	ledgerRel := filepath.Join(".stew", "ledgers", ledger)
	ledgerPath := filepath.Join(stewDir, "ledgers", ledger)
	ledgerInfo, err := os.Stat(ledgerPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return Result{}, fmt.Errorf("%w: ledger %q has no atomic entries directory; run `stew migrate atomic-entries` or `stew init` for a fresh repo", ErrMissingLedger, ledger)
		}
		return Result{}, fmt.Errorf("stat ledger storage: %w", err)
	}
	if !ledgerInfo.IsDir() {
		return Result{}, fmt.Errorf("%w: ledger %q has no atomic entries directory; run `stew migrate atomic-entries` or `stew init` for a fresh repo", ErrMissingLedger, ledger)
	}

	entryFiles, content, err := readEntryFiles(ledgerPath, ledgerRel)
	if err != nil {
		return Result{}, fmt.Errorf("read ledger %q: %w", ledger, err)
	}

	return Result{
		TargetDir:  targetDir,
		LedgerPath: ledgerRel,
		Content:    content,
		EntryFiles: entryFiles,
	}, nil
}

func validateLedgerName(name string) (string, error) {
	trimmed := strings.TrimSpace(name)
	if trimmed == "" {
		return "", fmt.Errorf("%w: ledger name is required", ErrInvalidLedger)
	}
	if trimmed == "stew" {
		return "", fmt.Errorf("%w: reserved ledger name %q", ErrInvalidLedger, trimmed)
	}
	if trimmed != name {
		return "", fmt.Errorf("%w: ledger name must not include surrounding whitespace", ErrInvalidLedger)
	}
	if trimmed == "." || trimmed == ".." || strings.Contains(trimmed, "..") {
		return "", fmt.Errorf("%w: ledger name must not contain path traversal", ErrInvalidLedger)
	}
	if strings.ContainsAny(trimmed, `/\`) {
		return "", fmt.Errorf("%w: ledger name must not contain path separators", ErrInvalidLedger)
	}
	if filepath.Clean(trimmed) != trimmed {
		return "", fmt.Errorf("%w: ledger name must be a base name", ErrInvalidLedger)
	}
	return trimmed, nil
}

func readEntryFiles(dir, ledgerRel string) ([]EntryFile, string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, "", err
	}

	names := make([]string, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".md") {
			continue
		}
		names = append(names, entry.Name())
	}
	sort.Strings(names)

	parts := make([]string, 0, len(names))
	entryFiles := make([]EntryFile, 0, len(names))
	for _, name := range names {
		path := filepath.Join(dir, name)
		bytes, err := os.ReadFile(path)
		if err != nil {
			return nil, "", err
		}
		content := strings.TrimRight(string(bytes), "\n") + "\n"
		entryFiles = append(entryFiles, EntryFile{
			Name:    name,
			RelPath: filepath.ToSlash(filepath.Join(ledgerRel, name)),
			Content: content,
		})
		parts = append(parts, content)
	}
	return entryFiles, strings.Join(parts, "\n"), nil
}
