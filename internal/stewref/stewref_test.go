package stewref

import (
	"errors"
	"testing"
)

func TestParseEntryRefRoundTripsCanonical(t *testing.T) {
	ref, err := Parse("entry:decisions/2026-06-20T191722Z-5cxsdb-use-ids-for-generated-entry-filenames.md")
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}
	if ref.Kind != KindEntry {
		t.Fatalf("Kind = %q, want %q", ref.Kind, KindEntry)
	}
	if ref.Payload != "decisions/2026-06-20T191722Z-5cxsdb-use-ids-for-generated-entry-filenames.md" {
		t.Fatalf("Payload = %q", ref.Payload)
	}
	if ref.String() != "entry:decisions/2026-06-20T191722Z-5cxsdb-use-ids-for-generated-entry-filenames.md" {
		t.Fatalf("String() = %q", ref.String())
	}
}

func TestParseAcceptsHistoricalEntryFilenameShape(t *testing.T) {
	ref, err := Parse("entry:iterations/2026-06-16T111551Z-implement-atomic-entry-storage.md")
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}
	if ref.String() != "entry:iterations/2026-06-16T111551Z-implement-atomic-entry-storage.md" {
		t.Fatalf("String() = %q", ref.String())
	}
}

func TestParseFileRefNormalizesRepoPath(t *testing.T) {
	ref, err := Parse(`file:./internal\stewentry/../stewref/stewref.go`)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}
	if ref.Kind != KindFile {
		t.Fatalf("Kind = %q, want %q", ref.Kind, KindFile)
	}
	if ref.String() != "file:internal/stewref/stewref.go" {
		t.Fatalf("String() = %q", ref.String())
	}
}

func TestEntryConstructorReturnsCanonicalRef(t *testing.T) {
	ref, err := Entry("decisions", "2026-06-20T191722Z-5cxsdb-use-ids-for-generated-entry-filenames.md")
	if err != nil {
		t.Fatalf("Entry() error = %v", err)
	}
	if ref.String() != "entry:decisions/2026-06-20T191722Z-5cxsdb-use-ids-for-generated-entry-filenames.md" {
		t.Fatalf("String() = %q", ref.String())
	}
}

func TestFileConstructorReturnsCanonicalRef(t *testing.T) {
	ref, err := File("./internal/stewref/stewref.go")
	if err != nil {
		t.Fatalf("File() error = %v", err)
	}
	if ref.String() != "file:internal/stewref/stewref.go" {
		t.Fatalf("String() = %q", ref.String())
	}
}

func TestParseRejectsInvalidRefs(t *testing.T) {
	cases := []string{
		"",
		"file",
		":internal/stewref/stewref.go",
		"file:",
		"commit:abc123",
		"file:/tmp/stew.go",
		`file:C:\tmp\stew.go`,
		"file:../outside.go",
		"file:internal/../../outside.go",
		"file: internal/stewref/stewref.go",
		"entry:decisions",
		"entry:/entry.md",
		"entry:decisions/",
		"entry:decisions/nested/entry.md",
		`entry:decisions\nested.md`,
		"entry:decisions/entry.txt",
		"entry:../decisions/entry.md",
		"entry:decisions/../entry.md",
		"entry: decisions/entry.md",
		"entry:decisions/ entry.md",
	}

	for _, value := range cases {
		t.Run(value, func(t *testing.T) {
			_, err := Parse(value)
			if !errors.Is(err, ErrInvalidRef) {
				t.Fatalf("Parse(%q) error = %v, want ErrInvalidRef", value, err)
			}
		})
	}
}
