package stewappend

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ankitvg/stew/internal/stewentry"
	"github.com/ankitvg/stew/internal/stewref"
)

var (
	ErrInvalidLedger  = errors.New("invalid ledger")
	ErrUnknownLedger  = errors.New("unknown ledger")
	ErrMissingBody    = errors.New("missing body source")
	ErrBodyConflict   = errors.New("multiple body sources")
	ErrMissingPrompt  = errors.New("missing prompt")
	ErrMissingSummary = errors.New("missing summary")
	ErrMissingStorage = errors.New("missing ledger storage")
)

type Options struct {
	TargetDir string
	Ledger    string
	Prompt    string
	Summary   string

	Message    string
	MessageSet bool
	FilePath   string
	Stdin      io.Reader
	StdinIsTTY func() bool

	Now        func() time.Time
	NewEntryID func() (string, error)
}

type Result struct {
	TargetDir  string
	EntryPath  string
	EntryRef   string
	LedgerPath string
}

func Run(opts Options) (Result, error) {
	if opts.TargetDir == "" {
		opts.TargetDir = "."
	}
	if opts.Now == nil {
		opts.Now = time.Now
	}
	if opts.NewEntryID == nil {
		opts.NewEntryID = stewentry.RandomID
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
	if strings.TrimSpace(opts.Prompt) == "" {
		return Result{}, ErrMissingPrompt
	}
	if strings.TrimSpace(opts.Summary) == "" {
		return Result{}, ErrMissingSummary
	}
	if strings.ContainsAny(opts.Summary, "\r\n") {
		return Result{}, fmt.Errorf("%w: summary must be a single line", ErrMissingSummary)
	}

	stewDir := filepath.Join(targetDir, ".stew")
	specPath := filepath.Join(stewDir, ledger+".spec.md")
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

	body, err := resolveBody(opts)
	if err != nil {
		return Result{}, err
	}

	entriesDir := filepath.Join(stewDir, "ledgers", ledger)
	entriesInfo, err := os.Stat(entriesDir)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return Result{}, fmt.Errorf("%w: ledger %q has no atomic entries directory; run `stew migrate atomic-entries` or `stew init` for a fresh repo", ErrMissingStorage, ledger)
		}
		return Result{}, fmt.Errorf("stat ledger storage: %w", err)
	}
	if !entriesInfo.IsDir() {
		return Result{}, fmt.Errorf("%w: ledger %q has no atomic entries directory; run `stew migrate atomic-entries` or `stew init` for a fresh repo", ErrMissingStorage, ledger)
	}

	now := opts.Now()
	timestamp := stewentry.FormatTimestamp(now)
	entry := stewentry.Render(now, opts.Summary, opts.Prompt, body)
	entryID, err := opts.NewEntryID()
	if err != nil {
		return Result{}, err
	}
	entryPath, err := writeEntryFile(entriesDir, timestamp, entryID, opts.Summary, entry)
	if err != nil {
		return Result{}, err
	}
	entryFileName := filepath.Base(entryPath)
	entryRef, err := stewref.Entry(ledger, entryFileName)
	if err != nil {
		return Result{}, fmt.Errorf("build entry ref: %w", err)
	}

	return Result{
		TargetDir:  targetDir,
		EntryPath:  filepath.Join(".stew", "ledgers", ledger, entryFileName),
		EntryRef:   entryRef.String(),
		LedgerPath: filepath.Join(".stew", "ledgers", ledger),
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

func resolveBody(opts Options) (string, error) {
	sourceCount := 0
	if opts.MessageSet {
		sourceCount++
	}
	if opts.FilePath != "" {
		sourceCount++
	}

	stdinAvailable := opts.Stdin != nil
	if stdinAvailable && opts.StdinIsTTY != nil && opts.StdinIsTTY() {
		stdinAvailable = false
	}
	if stdinAvailable {
		sourceCount++
	}

	if sourceCount == 0 {
		return "", fmt.Errorf("%w: provide body via stdin, -m, or -F", ErrMissingBody)
	}
	if sourceCount > 1 {
		return "", fmt.Errorf("%w: use exactly one of stdin, -m, or -F", ErrBodyConflict)
	}

	switch {
	case opts.MessageSet:
		return opts.Message, nil
	case opts.FilePath != "":
		bytes, err := os.ReadFile(opts.FilePath)
		if err != nil {
			return "", fmt.Errorf("read body file: %w", err)
		}
		return string(bytes), nil
	default:
		bytes, err := io.ReadAll(opts.Stdin)
		if err != nil {
			return "", fmt.Errorf("read stdin: %w", err)
		}
		return string(bytes), nil
	}
}

func writeEntryFile(entriesDir, timestamp, entryID, summary, entry string) (string, error) {
	for suffix := 1; ; suffix++ {
		name, err := stewentry.SuffixedFilenameWithID(timestamp, entryID, summary, suffix)
		if err != nil {
			return "", fmt.Errorf("build entry filename: %w", err)
		}
		path := filepath.Join(entriesDir, name)
		file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o644)
		if err != nil {
			if errors.Is(err, os.ErrExist) {
				continue
			}
			return "", fmt.Errorf("create entry file: %w", err)
		}
		if _, err := file.WriteString(entry); err != nil {
			_ = file.Close()
			return "", fmt.Errorf("write entry file: %w", err)
		}
		if err := file.Close(); err != nil {
			return "", fmt.Errorf("close entry file: %w", err)
		}
		return path, nil
	}
}
