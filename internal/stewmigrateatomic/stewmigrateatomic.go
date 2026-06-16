package stewmigrateatomic

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ankitvg/stew/internal/stewentry"
	"github.com/ankitvg/stew/internal/stewledgers"
)

type Options struct {
	TargetDir string
	DryRun    bool
}

type Result struct {
	TargetDir    string
	Ledgers      []LedgerResult
	LedgerCount  int
	EntryCount   int
	RemovedCount int
	SkippedCount int
	WouldRemove  int
	WouldWrite   int
	WouldMigrate int
}

type LedgerResult struct {
	Name       string
	LegacyPath string
	EntriesDir string
	Entries    int
	Skipped    bool
}

func Run(opts Options) (Result, error) {
	ledgersResult, err := stewledgers.List(stewledgers.Options{TargetDir: opts.TargetDir})
	if err != nil {
		return Result{}, err
	}

	result := Result{TargetDir: ledgersResult.TargetDir}
	prepared := make([]preparedLedger, 0, len(ledgersResult.Ledgers))

	for _, ledger := range ledgersResult.Ledgers {
		legacyPath := filepath.Join(ledgersResult.TargetDir, ".stew", ledger.Name+".md")
		legacyInfo, err := os.Stat(legacyPath)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				result.SkippedCount++
				result.Ledgers = append(result.Ledgers, LedgerResult{
					Name:       ledger.Name,
					LegacyPath: filepath.Join(".stew", ledger.Name+".md"),
					EntriesDir: filepath.Join(".stew", "ledgers", ledger.Name),
					Skipped:    true,
				})
				continue
			}
			return Result{}, fmt.Errorf("stat legacy ledger %q: %w", ledger.Name, err)
		}
		if legacyInfo.IsDir() {
			return Result{}, fmt.Errorf("legacy ledger %q is a directory", ledger.Name)
		}

		bytes, err := os.ReadFile(legacyPath)
		if err != nil {
			return Result{}, fmt.Errorf("read legacy ledger %q: %w", ledger.Name, err)
		}
		entries, err := parseLegacyLedger(ledger.Name, string(bytes))
		if err != nil {
			return Result{}, err
		}

		entryDir := filepath.Join(ledgersResult.TargetDir, ".stew", "ledgers", ledger.Name)
		prepared = append(prepared, preparedLedger{
			name:       ledger.Name,
			legacyPath: legacyPath,
			entryDir:   entryDir,
			entries:    entries,
		})
		result.LedgerCount++
		result.EntryCount += len(entries)
		result.Ledgers = append(result.Ledgers, LedgerResult{
			Name:       ledger.Name,
			LegacyPath: filepath.Join(".stew", ledger.Name+".md"),
			EntriesDir: filepath.Join(".stew", "ledgers", ledger.Name),
			Entries:    len(entries),
		})
	}

	if opts.DryRun {
		result.WouldMigrate = result.LedgerCount
		result.WouldWrite = result.EntryCount
		result.WouldRemove = result.LedgerCount
		return result, nil
	}

	written := make(map[string][]string, len(prepared))
	for _, ledger := range prepared {
		paths, err := writeLedgerEntries(ledger.entryDir, ledger.entries)
		if err != nil {
			return Result{}, fmt.Errorf("write atomic entries for %q: %w", ledger.name, err)
		}
		written[ledger.name] = paths
	}

	for _, ledger := range prepared {
		if err := verifyWrittenEntries(ledger.name, written[ledger.name]); err != nil {
			return Result{}, err
		}
	}

	for _, ledger := range prepared {
		if err := os.Remove(ledger.legacyPath); err != nil {
			return Result{}, fmt.Errorf("remove legacy ledger %q: %w", ledger.name, err)
		}
		result.RemovedCount++
	}

	return result, nil
}

type preparedLedger struct {
	name       string
	legacyPath string
	entryDir   string
	entries    []stewentry.Entry
}

func parseLegacyLedger(ledger, content string) ([]stewentry.Entry, error) {
	entries := stewentry.Parse(content)
	if len(entries) == 0 {
		if hasOnlyPreamble(content) {
			return nil, nil
		}
		return nil, fmt.Errorf("legacy ledger %q has no valid Stew entries", ledger)
	}
	return entries, nil
}

func hasOnlyPreamble(content string) bool {
	lines := strings.Split(strings.ReplaceAll(content, "\r\n", "\n"), "\n")
	remaining := make([]string, 0, len(lines))
	for _, raw := range lines {
		line := strings.TrimSpace(raw)
		if line == "" || strings.HasPrefix(line, "# ") || line == "<!-- Managed by stew -->" {
			continue
		}
		remaining = append(remaining, line)
	}
	return len(remaining) == 0
}

func writeLedgerEntries(entryDir string, entries []stewentry.Entry) ([]string, error) {
	if err := ensureWritableEntryDir(entryDir); err != nil {
		return nil, err
	}

	paths := make([]string, 0, len(entries))
	for _, entry := range entries {
		content := strings.TrimRight(entry.Content, "\n") + "\n"
		path, err := createEntryFile(entryDir, entry.Timestamp, entry.Summary, content)
		if err != nil {
			return nil, err
		}
		paths = append(paths, path)
	}
	return paths, nil
}

func ensureWritableEntryDir(entryDir string) error {
	info, err := os.Stat(entryDir)
	if err == nil {
		if !info.IsDir() {
			return fmt.Errorf("entry storage exists but is not a directory: %s", entryDir)
		}
		entries, err := os.ReadDir(entryDir)
		if err != nil {
			return err
		}
		for _, entry := range entries {
			if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".md") {
				return fmt.Errorf("entry storage already contains markdown entries: %s", entryDir)
			}
		}
		return nil
	}
	if !errors.Is(err, os.ErrNotExist) {
		return err
	}
	return os.MkdirAll(entryDir, 0o755)
}

func createEntryFile(entryDir, timestamp, summary, content string) (string, error) {
	for suffix := 1; ; suffix++ {
		name, err := stewentry.SuffixedFilename(timestamp, summary, suffix)
		if err != nil {
			return "", err
		}
		path := filepath.Join(entryDir, name)
		file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o644)
		if err != nil {
			if errors.Is(err, os.ErrExist) {
				continue
			}
			return "", err
		}
		if _, err := file.WriteString(content); err != nil {
			_ = file.Close()
			return "", err
		}
		if err := file.Close(); err != nil {
			return "", err
		}
		return path, nil
	}
}

func verifyWrittenEntries(ledger string, paths []string) error {
	for _, path := range paths {
		bytes, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("verify atomic entry for %q: %w", ledger, err)
		}
		if _, err := stewentry.ValidateContent(string(bytes)); err != nil {
			return fmt.Errorf("verify atomic entry for %q: %w", ledger, err)
		}
	}
	return nil
}
