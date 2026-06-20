package stewledgertail

import (
	"errors"
	"fmt"

	"github.com/ankitvg/stew/internal/stewentry"
	"github.com/ankitvg/stew/internal/stewledgercat"
	"github.com/ankitvg/stew/internal/stewref"
)

var ErrInvalidLimit = errors.New("invalid limit")

type Options struct {
	TargetDir string
	Ledger    string
	Limit     int
}

type Result struct {
	TargetDir  string
	LedgerPath string
	Content    string
	Entries    []Entry
}

type Entry struct {
	Ref       string
	Timestamp string
	Summary   string
	Prompt    string
	Body      string
}

type parsedEntry struct {
	entry   Entry
	content stewentry.Entry
}

func Run(opts Options) (Result, error) {
	if opts.Limit <= 0 {
		return Result{}, fmt.Errorf("%w: must be greater than zero", ErrInvalidLimit)
	}

	catResult, err := stewledgercat.Run(stewledgercat.Options{
		TargetDir: opts.TargetDir,
		Ledger:    opts.Ledger,
	})
	if err != nil {
		return Result{}, err
	}

	allEntries := make([]parsedEntry, 0, len(catResult.EntryFiles))
	for _, entryFile := range catResult.EntryFiles {
		entry, err := stewentry.ValidateContent(entryFile.Content)
		if err != nil {
			return Result{}, fmt.Errorf("parse entry file %s: %w", entryFile.RelPath, err)
		}
		ref, err := stewref.Entry(opts.Ledger, entryFile.Name)
		if err != nil {
			return Result{}, fmt.Errorf("build entry ref for %s: %w", entryFile.RelPath, err)
		}
		allEntries = append(allEntries, parsedEntry{
			entry: Entry{
				Ref:       ref.String(),
				Timestamp: entry.Timestamp,
				Summary:   entry.Summary,
				Prompt:    entry.Prompt,
				Body:      entry.Body,
			},
			content: entry,
		})
	}

	tail := tailParsedEntries(allEntries, opts.Limit)
	contentEntries := make([]stewentry.Entry, 0, len(tail))
	entries := make([]Entry, 0, len(tail))
	for _, entry := range tail {
		contentEntries = append(contentEntries, entry.content)
		entries = append(entries, entry.entry)
	}
	content := stewentry.Join(contentEntries)

	return Result{
		TargetDir:  catResult.TargetDir,
		LedgerPath: catResult.LedgerPath,
		Content:    content,
		Entries:    entries,
	}, nil
}

func tailParsedEntries(entries []parsedEntry, limit int) []parsedEntry {
	start := len(entries) - limit
	if start < 0 {
		start = 0
	}
	return entries[start:]
}
