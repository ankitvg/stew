package stewappend

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var (
	ErrInvalidLedger  = errors.New("invalid ledger")
	ErrUnknownLedger  = errors.New("unknown ledger")
	ErrMissingBody    = errors.New("missing body source")
	ErrBodyConflict   = errors.New("multiple body sources")
	ErrMissingPrompt  = errors.New("missing prompt")
	ErrMissingSummary = errors.New("missing summary")
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

	Now func() time.Time
}

type Result struct {
	TargetDir  string
	LedgerPath string
}

func Run(opts Options) (Result, error) {
	if opts.TargetDir == "" {
		opts.TargetDir = "."
	}
	if opts.Now == nil {
		opts.Now = time.Now
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

	ledgerPath := filepath.Join(stewDir, ledger+".md")
	entry := renderEntry(opts.Now(), opts.Summary, opts.Prompt, body)
	if err := appendEntry(ledgerPath, entry); err != nil {
		return Result{}, err
	}

	return Result{
		TargetDir:  targetDir,
		LedgerPath: filepath.Join(".stew", ledger+".md"),
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

func renderEntry(now time.Time, summary, prompt, body string) string {
	builder := &strings.Builder{}
	builder.WriteString("## ")
	builder.WriteString(formatTimestamp(now))
	builder.WriteString(" — ")
	builder.WriteString(summary)
	builder.WriteString("\n\n")

	if strings.ContainsAny(prompt, "\r\n") {
		builder.WriteString("**Prompt:**\n")
		builder.WriteString(strings.TrimRight(prompt, "\n"))
	} else {
		builder.WriteString("**Prompt:** ")
		builder.WriteString(prompt)
	}
	builder.WriteString("\n\n")

	builder.WriteString(strings.TrimRight(body, "\n"))
	builder.WriteString("\n\n---\n")
	return builder.String()
}

func appendEntry(path, entry string) error {
	existingBytes, err := os.ReadFile(path)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("read ledger: %w", err)
	}

	existing := strings.TrimRight(string(existingBytes), "\n")
	var next string
	if existing == "" {
		next = entry
	} else {
		next = existing + "\n\n" + entry
	}

	if err := os.WriteFile(path, []byte(next), 0o644); err != nil {
		return fmt.Errorf("write ledger: %w", err)
	}
	return nil
}

func formatTimestamp(now time.Time) string {
	return now.UTC().Format("2006-01-02T15:04:05Z")
}
