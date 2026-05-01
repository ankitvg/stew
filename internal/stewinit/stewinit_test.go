package stewinit

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestFormatTimestampUTCSecondPrecision(t *testing.T) {
	ts := time.Date(2026, 4, 26, 12, 5, 7, 987654321, time.FixedZone("UTC+1", 3600))
	got := formatTimestamp(ts)
	want := "2026-04-26T11:05:07Z"
	if got != want {
		t.Fatalf("formatTimestamp() = %q, want %q", got, want)
	}
}

func TestResolveUserNamePrecedence(t *testing.T) {
	t.Run("prefers git config", func(t *testing.T) {
		got := resolveUserName("/tmp", func(string) string { return "env-user" }, func(string) (string, error) {
			return "git-user", nil
		})
		if got != "git-user" {
			t.Fatalf("resolveUserName() = %q, want git-user", got)
		}
	})

	t.Run("falls back to USER", func(t *testing.T) {
		got := resolveUserName("/tmp", func(string) string { return "env-user" }, func(string) (string, error) {
			return "", errors.New("missing")
		})
		if got != "env-user" {
			t.Fatalf("resolveUserName() = %q, want env-user", got)
		}
	})

	t.Run("omits when unavailable", func(t *testing.T) {
		got := resolveUserName("/tmp", func(string) string { return "" }, func(string) (string, error) {
			return "", errors.New("missing")
		})
		if got != "" {
			t.Fatalf("resolveUserName() = %q, want empty string", got)
		}
	})
}

func TestUpsertManagedBlockInsertAndReplace(t *testing.T) {
	managed := renderManagedBlock()

	t.Run("inserts when missing", func(t *testing.T) {
		existing := "# Existing\n\nkeep this text\n"
		next, status, err := upsertManagedBlock(existing, managed)
		if err != nil {
			t.Fatalf("upsertManagedBlock() error = %v", err)
		}
		if status != FileStatusUpdated {
			t.Fatalf("status = %q, want %q", status, FileStatusUpdated)
		}
		if !strings.Contains(next, managedBlockStart) || !strings.Contains(next, managedBlockEnd) {
			t.Fatalf("managed block markers were not inserted")
		}
		if !strings.Contains(next, "keep this text") {
			t.Fatalf("existing content was not preserved")
		}
	})

	t.Run("replaces existing block only", func(t *testing.T) {
		existing := "intro\n" + managedBlockStart + "\nold\n" + managedBlockEnd + "\noutro\n"
		next, status, err := upsertManagedBlock(existing, managed)
		if err != nil {
			t.Fatalf("upsertManagedBlock() error = %v", err)
		}
		if status != FileStatusUpdated {
			t.Fatalf("status = %q, want %q", status, FileStatusUpdated)
		}
		if !strings.HasPrefix(next, "intro\n") || !strings.HasSuffix(next, "outro\n") {
			t.Fatalf("content outside managed block should be unchanged")
		}
		if strings.Count(next, managedBlockStart) != 1 || strings.Count(next, managedBlockEnd) != 1 {
			t.Fatalf("expected exactly one managed block")
		}
	})

	t.Run("errors on partial markers", func(t *testing.T) {
		existing := "intro\n" + managedBlockStart + "\nmissing end\n"
		_, _, err := upsertManagedBlock(existing, managed)
		if err == nil {
			t.Fatalf("expected error for partial managed block")
		}
		if !errors.Is(err, ErrCorruptManagedBlock) {
			t.Fatalf("error = %v, want ErrCorruptManagedBlock", err)
		}
	})
}

func TestRunFreshDirectoryCreatesStewFiles(t *testing.T) {
	tmp := t.TempDir()
	result, err := Run(Options{
		TargetDir: tmp,
		Now: func() time.Time {
			return time.Date(2026, 4, 26, 16, 42, 29, 0, time.UTC)
		},
		LookupEnv: func(string) string { return "env-user" },
		GitUserNameResolver: func(string) (string, error) {
			return "", errors.New("missing")
		},
		GitRepoChecker: func(string) error {
			return errors.New("not a repo")
		},
	})
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}

	if len(result.Warnings) == 0 {
		t.Fatalf("expected non-git warning")
	}

	required := []string{
		filepath.Join(".stew", "config.toml"),
		filepath.Join(".stew", "stew.spec.md"),
		filepath.Join(".stew", "iterations.md"),
		filepath.Join(".stew", "iterations.spec.md"),
		filepath.Join(".stew", "decisions.md"),
		filepath.Join(".stew", "decisions.spec.md"),
	}

	for _, rel := range required {
		if result.FileStatuses[rel] != FileStatusCreated {
			t.Fatalf("expected %s to be created, got %q", rel, result.FileStatuses[rel])
		}
		if _, err := os.Stat(filepath.Join(tmp, rel)); err != nil {
			t.Fatalf("expected %s to exist: %v", rel, err)
		}
	}

	if result.AgentsStatus != FileStatusCreated {
		t.Fatalf("AgentsStatus = %q, want %q", result.AgentsStatus, FileStatusCreated)
	}

	configBytes, err := os.ReadFile(filepath.Join(tmp, ".stew", "config.toml"))
	if err != nil {
		t.Fatalf("read config.toml: %v", err)
	}
	config := string(configBytes)
	if !strings.Contains(config, "created_at = \"2026-04-26T16:42:29Z\"") {
		t.Fatalf("config missing expected created_at: %s", config)
	}
	if !strings.Contains(config, "name = \"env-user\"") {
		t.Fatalf("config missing expected user name: %s", config)
	}

	agentsBytes, err := os.ReadFile(filepath.Join(tmp, "AGENTS.md"))
	if err != nil {
		t.Fatalf("read AGENTS.md: %v", err)
	}
	agents := string(agentsBytes)
	if !strings.Contains(agents, "durable project memory in append-only markdown ledgers") {
		t.Fatalf("AGENTS managed block missing durable memory description: %s", agents)
	}
	if !strings.Contains(agents, "Run `stew help` to discover available commands.") {
		t.Fatalf("AGENTS managed block missing help entrypoint: %s", agents)
	}
	if !strings.Contains(agents, "Run `stew full-spec` before working to load the workflow and ledger contract.") {
		t.Fatalf("AGENTS managed block missing full-spec entrypoint: %s", agents)
	}
	if strings.Contains(agents, "Default ledgers in this repo:") {
		t.Fatalf("AGENTS managed block should not list default ledgers: %s", agents)
	}
	if strings.Contains(agents, "stew ledger new <name>") {
		t.Fatalf("AGENTS managed block should not include detailed ledger creation guidance: %s", agents)
	}
	if strings.Contains(agents, "stew append <ledger>") {
		t.Fatalf("AGENTS managed block should not include detailed append guidance: %s", agents)
	}
	if strings.Contains(agents, "Iterations entry spec:") {
		t.Fatalf("AGENTS managed block should not list spec file paths directly: %s", agents)
	}
}

func TestStewSpecDocumentsCustomLedgerPrimitive(t *testing.T) {
	spec := renderStewSpec()
	required := []string{
		"A ledger has these durable properties:",
		"- Name: the command-facing identifier",
		"- Spec: a ledger-specific contract defines the ledger's purpose",
		"- Append-only semantics: never edit past entries.",
		"- CLI interface: use `stew ledgers` to discover ledger names, `stew ledger tail` or `stew ledger cat` to read entries, and `stew append` to write entries.",
		"- Chronological ordering: newest entries are appended at the bottom.",
		"- Entry boundaries: each entry starts with an H2 UTC ISO 8601 timestamp and summary.",
		"- Attribution: each entry includes `**Prompt:**` for the originating prompt.",
		"- Repo affiliation: ledgers belong to the repository where Stew is initialized.",
		"To add a custom ledger, run `stew ledger new <name>`.",
		"## Working With Stew",
		"Run `stew ledgers` to list writable ledger names and descriptions.",
		"Run `stew ledger tail --all --limit 5` to load recent project memory from every discovered writable ledger.",
		"At the end of meaningful work, append entries to the appropriate ledgers according to their specs.",
		"Use Stew commands as the interface for ledger access: read with `stew ledger tail` or `stew ledger cat`, and write with `stew append`.",
	}
	for _, want := range required {
		if !strings.Contains(spec, want) {
			t.Fatalf("stew spec missing %q: %s", want, spec)
		}
	}
}

func TestRunIsIdempotent(t *testing.T) {
	tmp := t.TempDir()
	opts := Options{
		TargetDir: tmp,
		Now: func() time.Time {
			return time.Date(2026, 4, 26, 16, 42, 29, 0, time.UTC)
		},
		LookupEnv: func(string) string { return "" },
		GitUserNameResolver: func(string) (string, error) {
			return "git-user", nil
		},
		GitRepoChecker: func(string) error { return nil },
	}

	_, err := Run(opts)
	if err != nil {
		t.Fatalf("first Run() error = %v", err)
	}

	result2, err := Run(opts)
	if err != nil {
		t.Fatalf("second Run() error = %v", err)
	}

	for rel, status := range result2.FileStatuses {
		if status != FileStatusExists {
			t.Fatalf("expected %s status to be exists, got %q", rel, status)
		}
	}

	if result2.AgentsStatus != FileStatusUnchanged {
		t.Fatalf("AgentsStatus = %q, want %q", result2.AgentsStatus, FileStatusUnchanged)
	}
}

func TestRunNoAgentsFlagSkipsAgentsFile(t *testing.T) {
	tmp := t.TempDir()
	result, err := Run(Options{
		TargetDir:  tmp,
		NoAgentsMD: true,
		GitRepoChecker: func(string) error {
			return nil
		},
	})
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if result.AgentsStatus != FileStatusSkipped {
		t.Fatalf("AgentsStatus = %q, want %q", result.AgentsStatus, FileStatusSkipped)
	}
	if _, err := os.Stat(filepath.Join(tmp, "AGENTS.md")); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("AGENTS.md should not exist, stat err = %v", err)
	}
}

func TestRunFailsOnCorruptAgentsMarkers(t *testing.T) {
	tmp := t.TempDir()
	corrupt := "header\n" + managedBlockStart + "\npartial\n"
	if err := os.WriteFile(filepath.Join(tmp, "AGENTS.md"), []byte(corrupt), 0o644); err != nil {
		t.Fatalf("write corrupt AGENTS.md: %v", err)
	}

	_, err := Run(Options{
		TargetDir: tmp,
		GitRepoChecker: func(string) error {
			return nil
		},
	})
	if err == nil {
		t.Fatalf("expected error for corrupt AGENTS.md")
	}
	if !errors.Is(err, ErrCorruptManagedBlock) {
		t.Fatalf("error = %v, want ErrCorruptManagedBlock", err)
	}
}
