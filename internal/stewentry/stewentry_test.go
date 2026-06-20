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

func TestFilenameWithIDUsesCompactTimestampIDAndCleanSlug(t *testing.T) {
	got, err := FilenameWithID("2026-06-16T17:30:12Z", "k7p3qx", " Add: Atomic Entry Storage! ")
	if err != nil {
		t.Fatalf("FilenameWithID() error = %v", err)
	}
	want := "2026-06-16T173012Z-k7p3qx-add-atomic-entry-storage.md"
	if got != want {
		t.Fatalf("FilenameWithID() = %q, want %q", got, want)
	}
}

func TestSuffixedFilenameWithIDKeepsSuffixAfterSlug(t *testing.T) {
	got, err := SuffixedFilenameWithID("2026-06-16T17:30:12Z", "k7p3qx", "Add atomic entry storage", 2)
	if err != nil {
		t.Fatalf("SuffixedFilenameWithID() error = %v", err)
	}
	want := "2026-06-16T173012Z-k7p3qx-add-atomic-entry-storage-2.md"
	if got != want {
		t.Fatalf("SuffixedFilenameWithID() = %q, want %q", got, want)
	}
}

func TestRandomIDUsesLowercaseBase32Shape(t *testing.T) {
	got, err := RandomID()
	if err != nil {
		t.Fatalf("RandomID() error = %v", err)
	}
	if err := ValidateID(got); err != nil {
		t.Fatalf("RandomID() = %q, invalid: %v", got, err)
	}
}

func TestValidateIDRejectsNonCanonicalIDs(t *testing.T) {
	cases := []string{"", "abcde", "abcdefg", "ABCDEF", "abc018", "abcde-"}
	for _, id := range cases {
		t.Run(id, func(t *testing.T) {
			if err := ValidateID(id); err == nil {
				t.Fatalf("ValidateID(%q) succeeded, want error", id)
			}
		})
	}
}

func TestSlugFallsBackForEmptyASCII(t *testing.T) {
	if got := Slug("!!!"); got != "entry" {
		t.Fatalf("Slug() = %q, want entry", got)
	}
}
