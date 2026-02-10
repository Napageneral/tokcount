package ignore

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadSpec_WithDefaultAndCustomPatterns(t *testing.T) {
	root := t.TempDir()

	if err := os.WriteFile(filepath.Join(root, ".gitignore"), []byte("ignored.txt\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, ".cartographerignore"), []byte("build/\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, ".tokcountignore"), []byte("custom.txt\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	spec, err := LoadSpec(root, ".tokcountignore")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !spec.MatchPath(filepath.Join(root, "ignored.txt"), false) {
		t.Fatalf("expected .gitignore pattern to match ignored.txt")
	}
	if !spec.MatchPath(filepath.Join(root, "build"), true) {
		t.Fatalf("expected .cartographerignore pattern to match build/")
	}
	if !spec.MatchPath(filepath.Join(root, "custom.txt"), false) {
		t.Fatalf("expected custom ignore pattern to match custom.txt")
	}
	if spec.MatchPath(filepath.Join(root, "kept.txt"), false) {
		t.Fatalf("did not expect kept.txt to be ignored")
	}
}
