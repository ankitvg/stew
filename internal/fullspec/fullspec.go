package fullspec

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

var (
	ErrMissingStewDir = errors.New("stew is not initialized")
	ErrNoSpecFiles    = errors.New("no stew ledger specs found")
)

type Result struct {
	TargetDir string
	SpecFiles []string
	Content   string
}

func Load(targetDir string) (Result, error) {
	if strings.TrimSpace(targetDir) == "" {
		targetDir = "."
	}

	absTarget, err := filepath.Abs(targetDir)
	if err != nil {
		return Result{}, fmt.Errorf("resolve target path: %w", err)
	}

	stewDir := filepath.Join(absTarget, ".stew")
	entries, err := os.ReadDir(stewDir)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return Result{}, fmt.Errorf("%w in %s", ErrMissingStewDir, absTarget)
		}
		return Result{}, fmt.Errorf("read stew metadata: %w", err)
	}

	specNames := make([]string, 0)
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if strings.HasSuffix(name, ".spec.md") {
			specNames = append(specNames, name)
		}
	}

	if len(specNames) == 0 {
		return Result{}, fmt.Errorf("%w in %s", ErrNoSpecFiles, absTarget)
	}

	sort.Slice(specNames, func(i, j int) bool {
		left := specNames[i]
		right := specNames[j]
		if left == "stew.spec.md" && right != "stew.spec.md" {
			return true
		}
		if right == "stew.spec.md" && left != "stew.spec.md" {
			return false
		}
		return left < right
	})

	files := make([]string, 0, len(specNames))
	sections := make([]string, 0, len(specNames))
	for _, name := range specNames {
		relPath := filepath.Join(".stew", name)
		absPath := filepath.Join(absTarget, relPath)
		bytes, err := os.ReadFile(absPath)
		if err != nil {
			return Result{}, fmt.Errorf("read %s: %w", specSourceLabel(name), err)
		}

		content := strings.TrimRight(string(bytes), "\n")
		sections = append(sections, fmt.Sprintf("<!-- %s -->\n%s", specSourceLabel(name), content))
		files = append(files, relPath)
	}

	result := Result{
		TargetDir: absTarget,
		SpecFiles: files,
		Content:   strings.Join(sections, "\n\n"),
	}

	if !strings.HasSuffix(result.Content, "\n") {
		result.Content += "\n"
	}

	return result, nil
}

func specSourceLabel(name string) string {
	if name == "stew.spec.md" {
		return "stew spec"
	}
	return strings.TrimSuffix(name, ".spec.md") + " spec"
}
