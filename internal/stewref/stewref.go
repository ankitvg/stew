package stewref

import (
	"errors"
	"fmt"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

const (
	KindEntry = "entry"
	KindFile  = "file"
)

var (
	ErrInvalidRef = errors.New("invalid ref")
)

var windowsDrivePattern = regexp.MustCompile(`^[A-Za-z]:`)

type Ref struct {
	Kind    string
	Payload string
}

func Parse(value string) (Ref, error) {
	kind, payload, found := strings.Cut(value, ":")
	if !found {
		return Ref{}, fmt.Errorf("%w: expected <kind>:<payload>", ErrInvalidRef)
	}
	if kind == "" {
		return Ref{}, fmt.Errorf("%w: kind is required", ErrInvalidRef)
	}
	if payload == "" {
		return Ref{}, fmt.Errorf("%w: payload is required", ErrInvalidRef)
	}

	switch kind {
	case KindEntry:
		ledger, fileName, err := parseEntryPayload(payload)
		if err != nil {
			return Ref{}, err
		}
		return Ref{Kind: KindEntry, Payload: ledger + "/" + fileName}, nil
	case KindFile:
		normalized, err := normalizeRepoPath(payload)
		if err != nil {
			return Ref{}, err
		}
		return Ref{Kind: KindFile, Payload: normalized}, nil
	default:
		return Ref{}, fmt.Errorf("%w: unsupported kind %q", ErrInvalidRef, kind)
	}
}

func (r Ref) String() string {
	if r.Kind == "" {
		return ""
	}
	return r.Kind + ":" + r.Payload
}

func Entry(ledger, fileName string) (Ref, error) {
	ledger, err := validateLedgerName(ledger)
	if err != nil {
		return Ref{}, err
	}
	fileName, err = validateEntryFileName(fileName)
	if err != nil {
		return Ref{}, err
	}
	return Ref{Kind: KindEntry, Payload: ledger + "/" + fileName}, nil
}

func File(repoPath string) (Ref, error) {
	normalized, err := normalizeRepoPath(repoPath)
	if err != nil {
		return Ref{}, err
	}
	return Ref{Kind: KindFile, Payload: normalized}, nil
}

func parseEntryPayload(payload string) (string, string, error) {
	ledger, fileName, found := strings.Cut(payload, "/")
	if !found {
		return "", "", fmt.Errorf("%w: entry payload must be <ledger>/<entry-file.md>", ErrInvalidRef)
	}
	if strings.Contains(fileName, "/") || strings.Contains(fileName, `\`) {
		return "", "", fmt.Errorf("%w: entry filename must be a base name", ErrInvalidRef)
	}
	ledger, err := validateLedgerName(ledger)
	if err != nil {
		return "", "", err
	}
	fileName, err = validateEntryFileName(fileName)
	if err != nil {
		return "", "", err
	}
	return ledger, fileName, nil
}

func validateLedgerName(name string) (string, error) {
	trimmed := strings.TrimSpace(name)
	if trimmed == "" {
		return "", fmt.Errorf("%w: ledger name is required", ErrInvalidRef)
	}
	if trimmed == "stew" {
		return "", fmt.Errorf("%w: reserved ledger name %q", ErrInvalidRef, trimmed)
	}
	if trimmed != name {
		return "", fmt.Errorf("%w: ledger name must not include surrounding whitespace", ErrInvalidRef)
	}
	if trimmed == "." || trimmed == ".." || strings.Contains(trimmed, "..") {
		return "", fmt.Errorf("%w: ledger name must not contain path traversal", ErrInvalidRef)
	}
	if strings.ContainsAny(trimmed, `/\`) {
		return "", fmt.Errorf("%w: ledger name must not contain path separators", ErrInvalidRef)
	}
	if filepath.Clean(trimmed) != trimmed {
		return "", fmt.Errorf("%w: ledger name must be a base name", ErrInvalidRef)
	}
	return trimmed, nil
}

func validateEntryFileName(name string) (string, error) {
	trimmed := strings.TrimSpace(name)
	if trimmed == "" {
		return "", fmt.Errorf("%w: entry filename is required", ErrInvalidRef)
	}
	if trimmed != name {
		return "", fmt.Errorf("%w: entry filename must not include surrounding whitespace", ErrInvalidRef)
	}
	if trimmed == "." || trimmed == ".." || strings.Contains(trimmed, "..") {
		return "", fmt.Errorf("%w: entry filename must not contain path traversal", ErrInvalidRef)
	}
	if strings.ContainsAny(trimmed, `/\`) {
		return "", fmt.Errorf("%w: entry filename must be a base name", ErrInvalidRef)
	}
	if !strings.HasSuffix(trimmed, ".md") {
		return "", fmt.Errorf("%w: entry filename must end with .md", ErrInvalidRef)
	}
	return trimmed, nil
}

func normalizeRepoPath(value string) (string, error) {
	if strings.TrimSpace(value) == "" {
		return "", fmt.Errorf("%w: file path is required", ErrInvalidRef)
	}
	if strings.TrimSpace(value) != value {
		return "", fmt.Errorf("%w: file path must not include surrounding whitespace", ErrInvalidRef)
	}
	if filepath.IsAbs(value) || strings.HasPrefix(value, "/") || windowsDrivePattern.MatchString(value) {
		return "", fmt.Errorf("%w: file path must be repo-relative", ErrInvalidRef)
	}

	slashPath := strings.ReplaceAll(value, `\`, "/")
	cleaned := path.Clean(slashPath)
	if cleaned == "." || cleaned == ".." || strings.HasPrefix(cleaned, "../") || strings.Contains(cleaned, "/../") {
		return "", fmt.Errorf("%w: file path must not contain path traversal", ErrInvalidRef)
	}
	return cleaned, nil
}
