package stewledgertail

import (
	"errors"
	"fmt"

	"github.com/ankitvg/stew/internal/stewentry"
	"github.com/ankitvg/stew/internal/stewledgercat"
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
	Timestamp string
	Summary   string
	Prompt    string
	Body      string
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

	allEntries := stewentry.Parse(catResult.Content)
	tail := stewentry.Tail(allEntries, opts.Limit)
	content := stewentry.Join(tail)
	entries := make([]Entry, 0, len(tail))
	for _, entry := range tail {
		entries = append(entries, Entry{
			Timestamp: entry.Timestamp,
			Summary:   entry.Summary,
			Prompt:    entry.Prompt,
			Body:      entry.Body,
		})
	}
	return Result{
		TargetDir:  catResult.TargetDir,
		LedgerPath: catResult.LedgerPath,
		Content:    content,
		Entries:    entries,
	}, nil
}
