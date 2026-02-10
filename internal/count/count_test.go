package count

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Napageneral/tokcount/internal/ignore"
	"github.com/Napageneral/tokcount/internal/tokenizer"
)

func TestRun_CountsAndIgnores(t *testing.T) {
	root := t.TempDir()

	if err := os.MkdirAll(filepath.Join(root, "src"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(root, "node_modules"), 0o755); err != nil {
		t.Fatal(err)
	}

	mainFile := []byte("package main\n\nfunc main() {}\n")
	if err := os.WriteFile(filepath.Join(root, "src", "main.go"), mainFile, 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "node_modules", "lib.js"), []byte("export const x = 1;\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "binary.bin"), []byte{0x00, 0x01, 0x02, 0x03}, 0o644); err != nil {
		t.Fatal(err)
	}

	ignoreSpec, err := ignore.LoadSpec(root, "")
	if err != nil {
		t.Fatal(err)
	}

	tok, err := tokenizer.New("estimate")
	if err != nil {
		t.Fatal(err)
	}

	result, err := Run(Options{
		Root:       root,
		Tokenizer:  tok,
		IgnoreSpec: ignoreSpec,
	})
	if err != nil {
		t.Fatal(err)
	}

	if result.TotalFiles != 1 {
		t.Fatalf("expected 1 scanned file, got %d", result.TotalFiles)
	}
	if result.IgnoredFiles != 2 {
		t.Fatalf("expected 2 ignored files, got %d", result.IgnoredFiles)
	}
	if result.TotalTokens <= 0 {
		t.Fatalf("expected positive token count, got %d", result.TotalTokens)
	}
	if result.DirectoryTokens["src"] <= 0 {
		t.Fatalf("expected src directory tokens to be > 0")
	}
}
