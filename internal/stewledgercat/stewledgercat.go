package stewledgercat

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
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
			return Result{}, fmt.Errorf("%w: %s not found", ErrUnknownLedger, specRel)
		}
		return Result{}, fmt.Errorf("stat ledger spec: %w", err)
	}
	if specInfo.IsDir() {
		return Result{}, fmt.Errorf("%w: %s is a directory", ErrUnknownLedger, specRel)
	}

	ledgerRel := filepath.Join(".stew", ledger+".md")
	ledgerPath := filepath.Join(stewDir, ledger+".md")
	ledgerInfo, err := os.Stat(ledgerPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return Result{}, fmt.Errorf("%w: %s not found", ErrMissingLedger, ledgerRel)
		}
		return Result{}, fmt.Errorf("stat ledger: %w", err)
	}
	if ledgerInfo.IsDir() {
		return Result{}, fmt.Errorf("%w: %s is a directory", ErrMissingLedger, ledgerRel)
	}

	bytes, err := os.ReadFile(ledgerPath)
	if err != nil {
		return Result{}, fmt.Errorf("read %s: %w", ledgerRel, err)
	}

	return Result{
		TargetDir:  targetDir,
		LedgerPath: ledgerRel,
		Content:    string(bytes),
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
