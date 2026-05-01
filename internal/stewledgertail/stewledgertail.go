package stewledgertail

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

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

var entryHeadingPattern = regexp.MustCompile(`(?m)^## (\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z) — (.+)$`)

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
	entries := parseEntries(content)
	return Result{
		TargetDir:  catResult.TargetDir,
		LedgerPath: catResult.LedgerPath,
		Content:    content,
		Entries:    entries,
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

func parseEntries(content string) []Entry {
	matches := entryHeadingPattern.FindAllStringSubmatchIndex(content, -1)
	entries := make([]Entry, 0, len(matches))
	for i, match := range matches {
		start := match[0]
		end := len(content)
		if i+1 < len(matches) {
			end = matches[i+1][0]
		}

		raw := content[start:end]
		timestamp := content[match[2]:match[3]]
		summary := content[match[4]:match[5]]
		prompt, body := parseEntryFields(raw)
		entries = append(entries, Entry{
			Timestamp: timestamp,
			Summary:   summary,
			Prompt:    prompt,
			Body:      body,
		})
	}
	return entries
}

func parseEntryFields(raw string) (string, string) {
	_, rest, found := strings.Cut(raw, "\n")
	if !found {
		return "", ""
	}
	rest = strings.TrimLeft(rest, "\n")
	if !strings.HasPrefix(rest, "**Prompt:**") {
		return "", stripEntryTerminator(rest)
	}

	afterPrompt := strings.TrimPrefix(rest, "**Prompt:**")
	afterPrompt = strings.TrimPrefix(afterPrompt, " ")
	afterPrompt = strings.TrimPrefix(afterPrompt, "\n")
	prompt, body, found := strings.Cut(afterPrompt, "\n\n")
	if !found {
		return strings.TrimRight(prompt, "\n"), ""
	}
	return strings.TrimRight(prompt, "\n"), stripEntryTerminator(body)
}

func stripEntryTerminator(body string) string {
	trimmed := strings.TrimRight(body, "\n")
	if trimmed == "---" {
		return ""
	}
	if strings.HasSuffix(trimmed, "\n---") {
		trimmed = strings.TrimSuffix(trimmed, "\n---")
	}
	return strings.TrimRight(trimmed, "\n")
}
