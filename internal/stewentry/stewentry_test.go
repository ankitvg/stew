package stewentry

import "testing"

func TestFilenameUsesCompactTimestampAndCleanSlug(t *testing.T) {
	got, err := Filename("2026-06-16T17:30:12Z", " Add: Atomic Entry Storage! ")
	if err != nil {
		t.Fatalf("Filename() error = %v", err)
	}
	want := "2026-06-16T173012Z-add-atomic-entry-storage.md"
	if got != want {
		t.Fatalf("Filename() = %q, want %q", got, want)
	}
}

func TestSlugFallsBackForEmptyASCII(t *testing.T) {
	if got := Slug("!!!"); got != "entry" {
		t.Fatalf("Slug() = %q, want entry", got)
	}
}
