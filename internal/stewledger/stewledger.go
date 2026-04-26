package stewledger

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var (
	ErrMissingStewDir = errors.New(".stew directory not found")
	ErrInvalidName    = errors.New("invalid ledger name")
	ErrLedgerExists   = errors.New("ledger already exists")
)

type Options struct {
	TargetDir   string
	Name        string
	Description string
	Threshold   string
}

type Result struct {
	TargetDir  string
	LedgerPath string
	SpecPath   string
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

	stewDir := filepath.Join(targetDir, ".stew")
	stewInfo, err := os.Stat(stewDir)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return Result{}, fmt.Errorf("%w in %s; run `stew init` first", ErrMissingStewDir, targetDir)
		}
		return Result{}, fmt.Errorf("stat .stew directory: %w", err)
	}
	if !stewInfo.IsDir() {
		return Result{}, fmt.Errorf("%w: %s is not a directory; run `stew init` first", ErrMissingStewDir, stewDir)
	}

	name, err := validateName(opts.Name)
	if err != nil {
		return Result{}, err
	}

	ledgerRel := filepath.Join(".stew", name+".md")
	specRel := filepath.Join(".stew", name+".spec.md")
	ledgerPath := filepath.Join(targetDir, ledgerRel)
	specPath := filepath.Join(targetDir, specRel)

	if err := ensureMissing(ledgerPath, ledgerRel); err != nil {
		return Result{}, err
	}
	if err := ensureMissing(specPath, specRel); err != nil {
		return Result{}, err
	}

	if err := os.WriteFile(ledgerPath, []byte(renderLedger(name)), 0o644); err != nil {
		return Result{}, fmt.Errorf("write %s: %w", ledgerRel, err)
	}
	if err := os.WriteFile(specPath, []byte(renderSpec(name, opts.Description, opts.Threshold)), 0o644); err != nil {
		return Result{}, fmt.Errorf("write %s: %w", specRel, err)
	}

	return Result{
		TargetDir:  targetDir,
		LedgerPath: ledgerRel,
		SpecPath:   specRel,
	}, nil
}

func validateName(name string) (string, error) {
	if name == "" {
		return "", fmt.Errorf("%w: name is required", ErrInvalidName)
	}
	if strings.TrimSpace(name) != name {
		return "", fmt.Errorf("%w: must not include surrounding whitespace", ErrInvalidName)
	}
	if strings.HasPrefix(name, ".") {
		return "", fmt.Errorf("%w: must not start with a dot", ErrInvalidName)
	}
	if len(name) > 40 {
		return "", fmt.Errorf("%w: must be 40 characters or fewer", ErrInvalidName)
	}
	if name[0] == '-' || name[len(name)-1] == '-' {
		return "", fmt.Errorf("%w: must not start or end with a hyphen", ErrInvalidName)
	}
	if reservedNames[name] {
		return "", fmt.Errorf("%w: reserved ledger name %q", ErrInvalidName, name)
	}

	previousHyphen := false
	for _, r := range name {
		switch {
		case r >= 'a' && r <= 'z':
			previousHyphen = false
		case r >= '0' && r <= '9':
			previousHyphen = false
		case r == '-':
			if previousHyphen {
				return "", fmt.Errorf("%w: must use single hyphens only", ErrInvalidName)
			}
			previousHyphen = true
		default:
			return "", fmt.Errorf("%w: use lowercase ASCII letters, digits, and single hyphens only", ErrInvalidName)
		}
	}

	return name, nil
}

var reservedNames = map[string]bool{
	"config":  true,
	"ledger":  true,
	"ledgers": true,
	"spec":    true,
	"specs":   true,
	"stew":    true,
	"stews":   true,
}

func ensureMissing(path, relPath string) error {
	_, err := os.Stat(path)
	if err == nil {
		return fmt.Errorf("%w: %s already exists", ErrLedgerExists, relPath)
	}
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	return fmt.Errorf("stat %s: %w", relPath, err)
}

func renderLedger(name string) string {
	return "# " + titleFromName(name) + "\n\n" +
		"<!-- Managed by stew -->\n"
}

func renderSpec(name, description, threshold string) string {
	description = strings.TrimSpace(description)
	threshold = strings.TrimSpace(threshold)
	if description == "" {
		description = "TODO: Describe what this ledger records and who should read it."
	}
	if threshold == "" {
		threshold = "TODO: Explain when an entry belongs in this ledger instead of another ledger."
	}

	return "# " + titleFromName(name) + " Spec\n\n" +
		description + "\n\n" +
		"## Body\n\n" +
		"Entries use freeform markdown under the standard Stew entry header. Keep each entry focused on one durable observation, change, or decision relevant to this ledger.\n\n" +
		"## When to append\n\n" +
		threshold + "\n"
}

func titleFromName(name string) string {
	parts := strings.Split(name, "-")
	for i, part := range parts {
		if part == "" {
			continue
		}
		parts[i] = strings.ToUpper(part[:1]) + part[1:]
	}
	return strings.Join(parts, " ")
}
