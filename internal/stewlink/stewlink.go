package stewlink

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/ankitvg/stew/internal/stewentry"
	"github.com/ankitvg/stew/internal/stewref"
)

const Version = 1

var ErrMalformedLink = errors.New("malformed link")

type Link struct {
	Version   int    `json:"version"`
	CreatedAt string `json:"createdAt"`
	Source    string `json:"source"`
	Target    string `json:"target"`
}

type CreateOptions struct {
	TargetDir string
	Source    stewref.Ref
	Target    stewref.Ref

	Now   func() time.Time
	NewID func() (string, error)
}

type CreateResult struct {
	TargetDir string
	LinkPath  string
	Link      Link
}

type ListOptions struct {
	TargetDir string
	Ref       stewref.Ref
}

type ListResult struct {
	TargetDir string
	Ref       string
	Links     []Link
}

func Create(opts CreateOptions) (CreateResult, error) {
	if opts.TargetDir == "" {
		opts.TargetDir = "."
	}
	if opts.Now == nil {
		opts.Now = time.Now
	}
	if opts.NewID == nil {
		opts.NewID = stewentry.RandomID
	}

	targetDir, err := resolveTargetDir(opts.TargetDir)
	if err != nil {
		return CreateResult{}, err
	}

	source, err := stewref.Resolve(stewref.ResolveOptions{
		TargetDir:     targetDir,
		Ref:           opts.Source,
		RequireExists: true,
	})
	if err != nil {
		return CreateResult{}, fmt.Errorf("resolve source ref: %w", err)
	}
	target, err := stewref.Resolve(stewref.ResolveOptions{
		TargetDir:     targetDir,
		Ref:           opts.Target,
		RequireExists: true,
	})
	if err != nil {
		return CreateResult{}, fmt.Errorf("resolve target ref: %w", err)
	}

	createdAt := stewentry.FormatTimestamp(opts.Now())
	link := Link{
		Version:   Version,
		CreatedAt: createdAt,
		Source:    source.Ref.String(),
		Target:    target.Ref.String(),
	}
	linkID, err := opts.NewID()
	if err != nil {
		return CreateResult{}, err
	}
	if err := stewentry.ValidateID(linkID); err != nil {
		return CreateResult{}, err
	}

	linksDir := filepath.Join(targetDir, ".stew", "links")
	if err := os.MkdirAll(linksDir, 0o755); err != nil {
		return CreateResult{}, fmt.Errorf("create links storage: %w", err)
	}
	path, err := writeLinkFile(linksDir, createdAt, linkID, link)
	if err != nil {
		return CreateResult{}, err
	}

	return CreateResult{
		TargetDir: targetDir,
		LinkPath:  filepath.ToSlash(filepath.Join(".stew", "links", filepath.Base(path))),
		Link:      link,
	}, nil
}

func List(opts ListOptions) (ListResult, error) {
	if opts.TargetDir == "" {
		opts.TargetDir = "."
	}

	targetDir, err := resolveTargetDir(opts.TargetDir)
	if err != nil {
		return ListResult{}, err
	}
	ref, err := stewref.Parse(opts.Ref.String())
	if err != nil {
		return ListResult{}, err
	}

	all, err := readLinks(filepath.Join(targetDir, ".stew", "links"))
	if err != nil {
		return ListResult{}, err
	}
	refString := ref.String()
	links := make([]Link, 0, len(all))
	for _, link := range all {
		if link.Source == refString || link.Target == refString {
			links = append(links, link)
		}
	}
	return ListResult{
		TargetDir: targetDir,
		Ref:       refString,
		Links:     links,
	}, nil
}

func resolveTargetDir(targetPath string) (string, error) {
	targetDir, err := filepath.Abs(targetPath)
	if err != nil {
		return "", fmt.Errorf("resolve target path: %w", err)
	}
	info, err := os.Stat(targetDir)
	if err != nil {
		return "", fmt.Errorf("stat target path: %w", err)
	}
	if !info.IsDir() {
		return "", fmt.Errorf("target path is not a directory: %s", targetDir)
	}
	return targetDir, nil
}

func writeLinkFile(dir, createdAt, linkID string, link Link) (string, error) {
	compact, err := stewentry.CompactTimestamp(createdAt)
	if err != nil {
		return "", err
	}
	body, err := json.MarshalIndent(link, "", "  ")
	if err != nil {
		return "", fmt.Errorf("encode link: %w", err)
	}
	body = append(body, '\n')

	for suffix := 1; ; suffix++ {
		name := compact + "-" + linkID
		if suffix > 1 {
			name = fmt.Sprintf("%s-%d", name, suffix)
		}
		path := filepath.Join(dir, name+".json")
		file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o644)
		if err != nil {
			if errors.Is(err, os.ErrExist) {
				continue
			}
			return "", fmt.Errorf("create link file: %w", err)
		}
		if _, err := file.Write(body); err != nil {
			_ = file.Close()
			return "", fmt.Errorf("write link file: %w", err)
		}
		if err := file.Close(); err != nil {
			return "", fmt.Errorf("close link file: %w", err)
		}
		return path, nil
	}
}

func readLinks(dir string) ([]Link, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return []Link{}, nil
		}
		return nil, err
	}

	names := make([]string, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".json" {
			continue
		}
		names = append(names, entry.Name())
	}
	sort.Strings(names)

	links := make([]Link, 0, len(names))
	for _, name := range names {
		link, err := readLinkFile(filepath.Join(dir, name))
		if err != nil {
			return nil, err
		}
		links = append(links, link)
	}
	return links, nil
}

func readLinkFile(path string) (Link, error) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return Link{}, err
	}
	var link Link
	if err := json.Unmarshal(bytes, &link); err != nil {
		return Link{}, fmt.Errorf("%w: %s: %v", ErrMalformedLink, path, err)
	}
	if err := validateStoredLink(link); err != nil {
		return Link{}, fmt.Errorf("%w: %s: %v", ErrMalformedLink, path, err)
	}
	return link, nil
}

func validateStoredLink(link Link) error {
	if link.Version != Version {
		return fmt.Errorf("version must be %d", Version)
	}
	if _, err := time.Parse("2006-01-02T15:04:05Z", link.CreatedAt); err != nil {
		return fmt.Errorf("createdAt must be UTC second timestamp")
	}
	source, err := stewref.Parse(link.Source)
	if err != nil {
		return fmt.Errorf("source: %w", err)
	}
	target, err := stewref.Parse(link.Target)
	if err != nil {
		return fmt.Errorf("target: %w", err)
	}
	if link.Source != source.String() {
		return fmt.Errorf("source must be canonical")
	}
	if link.Target != target.String() {
		return fmt.Errorf("target must be canonical")
	}
	return nil
}
