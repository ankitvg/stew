package stewledgertail

import (
	"errors"
	"fmt"
	"regexp"

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
}

var entryHeadingPattern = regexp.MustCompile(`(?m)^## \d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z — .+$`)

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

	content := tailEntries(catResult.Content, opts.Limit)
	return Result{
		TargetDir:  catResult.TargetDir,
		LedgerPath: catResult.LedgerPath,
		Content:    content,
	}, nil
}

func tailEntries(content string, limit int) string {
	matches := entryHeadingPattern.FindAllStringIndex(content, -1)
	if len(matches) == 0 {
		return ""
	}

	startIndex := len(matches) - limit
	if startIndex < 0 {
		startIndex = 0
	}
	return content[matches[startIndex][0]:]
}
