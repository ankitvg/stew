package stewref

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var (
	ErrMissingTarget = errors.New("missing ref target")
	ErrUnknownLedger = errors.New("unknown ledger")
)

type ResolveOptions struct {
	TargetDir     string
	Ref           Ref
	RequireExists bool
}

type ResolvedRef struct {
	Ref     Ref
	AbsPath string
	RelPath string
}

func Resolve(opts ResolveOptions) (ResolvedRef, error) {
	if opts.TargetDir == "" {
		opts.TargetDir = "."
	}

	targetDir, err := filepath.Abs(opts.TargetDir)
	if err != nil {
		return ResolvedRef{}, fmt.Errorf("resolve target path: %w", err)
	}
	info, err := os.Stat(targetDir)
	if err != nil {
		return ResolvedRef{}, fmt.Errorf("stat target path: %w", err)
	}
	if !info.IsDir() {
		return ResolvedRef{}, fmt.Errorf("target path is not a directory: %s", targetDir)
	}

	ref, err := Parse(opts.Ref.String())
	if err != nil {
		return ResolvedRef{}, err
	}

	relPath, err := resolveRelPath(targetDir, ref)
	if err != nil {
		return ResolvedRef{}, err
	}
	absPath := filepath.Join(targetDir, filepath.FromSlash(relPath))
	if opts.RequireExists {
		if err := requireFile(absPath); err != nil {
			return ResolvedRef{}, err
		}
	}

	return ResolvedRef{
		Ref:     ref,
		AbsPath: absPath,
		RelPath: relPath,
	}, nil
}

func resolveRelPath(targetDir string, ref Ref) (string, error) {
	switch ref.Kind {
	case KindFile:
		if ref.Payload == ".stew/ledgers" || strings.HasPrefix(ref.Payload, ".stew/ledgers/") {
			return "", fmt.Errorf("%w: ledger entries must use entry refs", ErrInvalidRef)
		}
		return ref.Payload, nil
	case KindEntry:
		ledger, fileName, err := parseEntryPayload(ref.Payload)
		if err != nil {
			return "", err
		}
		specPath := filepath.Join(targetDir, ".stew", ledger+".spec.md")
		specInfo, err := os.Stat(specPath)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				return "", fmt.Errorf("%w: ledger %q is not defined", ErrUnknownLedger, ledger)
			}
			return "", fmt.Errorf("stat ledger spec: %w", err)
		}
		if specInfo.IsDir() {
			return "", fmt.Errorf("%w: ledger %q is not defined", ErrUnknownLedger, ledger)
		}
		return ".stew/ledgers/" + ledger + "/" + fileName, nil
	default:
		return "", fmt.Errorf("%w: unsupported kind %q", ErrInvalidRef, ref.Kind)
	}
}

func requireFile(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("%w: %s", ErrMissingTarget, path)
		}
		return fmt.Errorf("stat ref target: %w", err)
	}
	if info.IsDir() {
		return fmt.Errorf("%w: target is a directory: %s", ErrMissingTarget, path)
	}
	return nil
}
