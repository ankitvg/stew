package stewinit

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type FileStatus string

const (
	FileStatusCreated   FileStatus = "created"
	FileStatusExists    FileStatus = "exists"
	FileStatusUpdated   FileStatus = "updated"
	FileStatusUnchanged FileStatus = "unchanged"
	FileStatusSkipped   FileStatus = "skipped"
)

var ErrCorruptManagedBlock = errors.New("AGENTS.md contains a partial or malformed Stew managed block")

type Options struct {
	TargetDir  string
	NoAgentsMD bool

	Now                 func() time.Time
	LookupEnv           func(string) string
	GitUserNameResolver func(string) (string, error)
	GitRepoChecker      func(string) error
}

type Result struct {
	TargetDir    string
	FileStatuses map[string]FileStatus
	AgentsStatus FileStatus
	Warnings     []string
}

func Run(opts Options) (Result, error) {
	if opts.TargetDir == "" {
		opts.TargetDir = "."
	}
	if opts.Now == nil {
		opts.Now = time.Now
	}
	if opts.LookupEnv == nil {
		opts.LookupEnv = os.Getenv
	}
	if opts.GitUserNameResolver == nil {
		opts.GitUserNameResolver = defaultGitUserNameResolver
	}
	if opts.GitRepoChecker == nil {
		opts.GitRepoChecker = defaultGitRepoChecker
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

	result := Result{
		TargetDir:    targetDir,
		FileStatuses: make(map[string]FileStatus),
	}

	if err := opts.GitRepoChecker(targetDir); err != nil {
		result.Warnings = append(result.Warnings, "target path is not a git repository; continuing")
	}

	username := resolveUserName(targetDir, opts.LookupEnv, opts.GitUserNameResolver)
	createdAt := formatTimestamp(opts.Now())

	stewDir := filepath.Join(targetDir, ".stew")
	if err := os.MkdirAll(stewDir, 0o755); err != nil {
		return Result{}, fmt.Errorf("create stew metadata: %w", err)
	}

	templates := []struct {
		relPath string
		body    string
	}{
		{relPath: filepath.Join(".stew", "config.toml"), body: renderConfig(createdAt, username)},
		{relPath: filepath.Join(".stew", "stew.spec.md"), body: renderStewSpec()},
		{relPath: filepath.Join(".stew", "iterations.spec.md"), body: renderIterationsSpec()},
		{relPath: filepath.Join(".stew", "decisions.spec.md"), body: renderDecisionsSpec()},
	}

	for _, file := range templates {
		status, err := createIfMissing(filepath.Join(targetDir, file.relPath), file.body)
		if err != nil {
			return Result{}, err
		}
		result.FileStatuses[file.relPath] = status
	}

	for _, ledger := range []string{"iterations", "decisions"} {
		relPath := filepath.Join(".stew", "ledgers", ledger)
		status, err := createLedgerStorageIfMissing(targetDir, ledger)
		if err != nil {
			return Result{}, err
		}
		result.FileStatuses[relPath] = status
	}

	linksRelPath := filepath.Join(".stew", "links")
	linksStatus, err := createLinksStorageIfMissing(targetDir)
	if err != nil {
		return Result{}, err
	}
	result.FileStatuses[linksRelPath] = linksStatus

	if opts.NoAgentsMD {
		result.AgentsStatus = FileStatusSkipped
		return result, nil
	}

	agentsPath := filepath.Join(targetDir, "AGENTS.md")
	agentsStatus, err := ensureAgentsManagedBlock(agentsPath, renderManagedBlock())
	if err != nil {
		return Result{}, err
	}
	result.AgentsStatus = agentsStatus

	return result, nil
}

func createLedgerStorageIfMissing(targetDir, ledger string) (FileStatus, error) {
	dirPath := filepath.Join(targetDir, ".stew", "ledgers", ledger)
	info, err := os.Stat(dirPath)
	if err == nil {
		if !info.IsDir() {
			return "", fmt.Errorf("ledger storage is not a directory: %s", dirPath)
		}
		return FileStatusExists, nil
	}
	if !errors.Is(err, os.ErrNotExist) {
		return "", fmt.Errorf("stat ledger storage: %w", err)
	}

	legacyPath := filepath.Join(targetDir, ".stew", ledger+".md")
	if legacyInfo, err := os.Stat(legacyPath); err == nil {
		if !legacyInfo.IsDir() {
			return FileStatusSkipped, nil
		}
	} else if !errors.Is(err, os.ErrNotExist) {
		return "", fmt.Errorf("stat legacy ledger: %w", err)
	}

	if err := os.MkdirAll(dirPath, 0o755); err != nil {
		return "", fmt.Errorf("create ledger storage: %w", err)
	}
	return FileStatusCreated, nil
}

func createLinksStorageIfMissing(targetDir string) (FileStatus, error) {
	dirPath := filepath.Join(targetDir, ".stew", "links")
	info, err := os.Stat(dirPath)
	if err == nil {
		if !info.IsDir() {
			return "", fmt.Errorf("links storage is not a directory: %s", dirPath)
		}
		return FileStatusExists, nil
	}
	if !errors.Is(err, os.ErrNotExist) {
		return "", fmt.Errorf("stat links storage: %w", err)
	}
	if err := os.MkdirAll(dirPath, 0o755); err != nil {
		return "", fmt.Errorf("create links storage: %w", err)
	}
	return FileStatusCreated, nil
}

func createIfMissing(path, body string) (FileStatus, error) {
	if _, err := os.Stat(path); err == nil {
		return FileStatusExists, nil
	} else if !errors.Is(err, os.ErrNotExist) {
		return "", fmt.Errorf("stat stew-managed file: %w", err)
	}

	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		return "", fmt.Errorf("write stew-managed file: %w", err)
	}
	return FileStatusCreated, nil
}

func ensureAgentsManagedBlock(path string, managedBlock string) (FileStatus, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			if err := os.WriteFile(path, []byte(managedBlock), 0o644); err != nil {
				return "", fmt.Errorf("write AGENTS.md: %w", err)
			}
			return FileStatusCreated, nil
		}
		return "", fmt.Errorf("read AGENTS.md: %w", err)
	}

	nextContent, status, err := upsertManagedBlock(string(content), managedBlock)
	if err != nil {
		return "", err
	}
	if status == FileStatusUnchanged {
		return status, nil
	}
	if err := os.WriteFile(path, []byte(nextContent), 0o644); err != nil {
		return "", fmt.Errorf("write AGENTS.md: %w", err)
	}
	return status, nil
}

func upsertManagedBlock(existing, managedBlock string) (string, FileStatus, error) {
	beginPresent := strings.Contains(existing, managedBlockStart)
	endPresent := strings.Contains(existing, managedBlockEnd)
	if beginPresent != endPresent {
		return "", "", fmt.Errorf("%w: repair markers %q and %q", ErrCorruptManagedBlock, managedBlockStart, managedBlockEnd)
	}

	if !beginPresent && !endPresent {
		trimmed := strings.TrimRight(existing, "\n")
		if trimmed == "" {
			return managedBlock, FileStatusUpdated, nil
		}
		return trimmed + "\n\n" + managedBlock, FileStatusUpdated, nil
	}

	beginIndex := strings.Index(existing, managedBlockStart)
	endIndex := strings.Index(existing, managedBlockEnd)
	if beginIndex < 0 || endIndex < 0 || beginIndex > endIndex {
		return "", "", fmt.Errorf("%w: malformed marker ordering", ErrCorruptManagedBlock)
	}
	endIndex += len(managedBlockEnd)
	if strings.HasPrefix(existing[endIndex:], "\r\n") {
		endIndex += 2
	} else if strings.HasPrefix(existing[endIndex:], "\n") {
		endIndex += 1
	}
	next := existing[:beginIndex] + managedBlock + existing[endIndex:]
	if next == existing {
		return existing, FileStatusUnchanged, nil
	}
	return next, FileStatusUpdated, nil
}

func formatTimestamp(now time.Time) string {
	return now.UTC().Format("2006-01-02T15:04:05Z")
}

func resolveUserName(targetDir string, lookupEnv func(string) string, gitNameResolver func(string) (string, error)) string {
	if gitNameResolver != nil {
		if name, err := gitNameResolver(targetDir); err == nil {
			if trimmed := strings.TrimSpace(name); trimmed != "" {
				return trimmed
			}
		}
	}
	if lookupEnv != nil {
		if envUser := strings.TrimSpace(lookupEnv("USER")); envUser != "" {
			return envUser
		}
	}
	return ""
}

func renderConfig(createdAt, userName string) string {
	builder := &strings.Builder{}
	builder.WriteString("[stew]\n")
	builder.WriteString("version = \"1\"\n")
	builder.WriteString(fmt.Sprintf("created_at = \"%s\"\n", createdAt))

	if strings.TrimSpace(userName) != "" {
		builder.WriteString("\n[user]\n")
		builder.WriteString(fmt.Sprintf("name = \"%s\"\n", escapeTOMLString(strings.TrimSpace(userName))))
	}

	return builder.String()
}

func escapeTOMLString(value string) string {
	value = strings.ReplaceAll(value, "\\", "\\\\")
	value = strings.ReplaceAll(value, "\"", "\\\"")
	return value
}

func defaultGitUserNameResolver(targetDir string) (string, error) {
	cmd := exec.Command("git", "-C", targetDir, "config", "user.name")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

func defaultGitRepoChecker(targetDir string) error {
	cmd := exec.Command("git", "-C", targetDir, "rev-parse", "--is-inside-work-tree")
	output, err := cmd.Output()
	if err != nil {
		return err
	}
	if strings.TrimSpace(string(output)) != "true" {
		return errors.New("not inside git work tree")
	}
	return nil
}
