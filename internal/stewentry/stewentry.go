package stewentry

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Entry struct {
	Timestamp string
	Summary   string
	Prompt    string
	Body      string
	Content   string
}

var HeadingPattern = regexp.MustCompile(`(?m)^## (\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z) — (.+)$`)

func Render(now time.Time, summary, prompt, body string) string {
	builder := &strings.Builder{}
	builder.WriteString("## ")
	builder.WriteString(FormatTimestamp(now))
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

func FormatTimestamp(now time.Time) string {
	return now.UTC().Format("2006-01-02T15:04:05Z")
}

func CompactTimestamp(timestamp string) (string, error) {
	parsed, err := time.Parse("2006-01-02T15:04:05Z", timestamp)
	if err != nil {
		return "", err
	}
	return parsed.UTC().Format("2006-01-02T150405Z"), nil
}

func Filename(timestamp, summary string) (string, error) {
	compact, err := CompactTimestamp(timestamp)
	if err != nil {
		return "", err
	}
	return compact + "-" + Slug(summary) + ".md", nil
}

func SuffixedFilename(timestamp, summary string, suffix int) (string, error) {
	if suffix <= 1 {
		return Filename(timestamp, summary)
	}
	compact, err := CompactTimestamp(timestamp)
	if err != nil {
		return "", err
	}
	return compact + "-" + Slug(summary) + "-" + strconv.Itoa(suffix) + ".md", nil
}

func Slug(value string) string {
	lowered := strings.ToLower(value)
	var builder strings.Builder
	previousHyphen := false
	for _, r := range lowered {
		switch {
		case r >= 'a' && r <= 'z':
			builder.WriteRune(r)
			previousHyphen = false
		case r >= '0' && r <= '9':
			builder.WriteRune(r)
			previousHyphen = false
		default:
			if !previousHyphen {
				builder.WriteByte('-')
				previousHyphen = true
			}
		}
	}

	slug := strings.Trim(builder.String(), "-")
	if slug == "" {
		return "entry"
	}
	if len(slug) > 80 {
		slug = strings.TrimRight(slug[:80], "-")
	}
	if slug == "" {
		return "entry"
	}
	return slug
}

func Parse(content string) []Entry {
	matches := HeadingPattern.FindAllStringSubmatchIndex(content, -1)
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
		prompt, body := parseFields(raw)
		entries = append(entries, Entry{
			Timestamp: timestamp,
			Summary:   summary,
			Prompt:    prompt,
			Body:      body,
			Content:   raw,
		})
	}
	return entries
}

func Tail(entries []Entry, limit int) []Entry {
	start := len(entries) - limit
	if start < 0 {
		start = 0
	}
	return entries[start:]
}

func Join(entries []Entry) string {
	parts := make([]string, 0, len(entries))
	for _, entry := range entries {
		parts = append(parts, strings.TrimRight(entry.Content, "\n")+"\n")
	}
	return strings.Join(parts, "\n")
}

func ValidateContent(content string) (Entry, error) {
	entries := Parse(content)
	if len(entries) != 1 {
		return Entry{}, fmt.Errorf("entry file must contain exactly one entry, found %d", len(entries))
	}
	if strings.TrimSpace(entries[0].Content) != strings.TrimSpace(content) {
		return Entry{}, fmt.Errorf("entry file contains content outside the entry")
	}
	return entries[0], nil
}

func parseFields(raw string) (string, string) {
	_, rest, found := strings.Cut(raw, "\n")
	if !found {
		return "", ""
	}
	rest = strings.TrimLeft(rest, "\n")
	if !strings.HasPrefix(rest, "**Prompt:**") {
		return "", StripTerminator(rest)
	}

	afterPrompt := strings.TrimPrefix(rest, "**Prompt:**")
	afterPrompt = strings.TrimPrefix(afterPrompt, " ")
	afterPrompt = strings.TrimPrefix(afterPrompt, "\n")
	prompt, body, found := strings.Cut(afterPrompt, "\n\n")
	if !found {
		return strings.TrimRight(prompt, "\n"), ""
	}
	return strings.TrimRight(prompt, "\n"), StripTerminator(body)
}

func StripTerminator(body string) string {
	trimmed := strings.TrimRight(body, "\n")
	if trimmed == "---" {
		return ""
	}
	if strings.HasSuffix(trimmed, "\n---") {
		trimmed = strings.TrimSuffix(trimmed, "\n---")
	}
	return strings.TrimRight(trimmed, "\n")
}
