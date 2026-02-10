package count

import (
	"bytes"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/Napageneral/tokcount/internal/ignore"
	"github.com/Napageneral/tokcount/internal/tokenizer"
)

const defaultMaxFileBytes int64 = 10 * 1024 * 1024 // 10 MiB

// Options controls repository counting behavior.
type Options struct {
	Root         string
	Tokenizer    tokenizer.Tokenizer
	IgnoreSpec   *ignore.Spec
	MaxFileBytes int64
}

// Result is the normalized token counting output.
type Result struct {
	Repository      string         `json:"repository"`
	Tokenizer       string         `json:"tokenizer"`
	TokenizerDetail string         `json:"-"`
	TotalTokens     int            `json:"total_tokens"`
	TotalFiles      int            `json:"total_files"`
	IgnoredFiles    int            `json:"ignored_files"`
	TotalLines      int            `json:"total_lines"`
	DirectoryTokens map[string]int `json:"-"`
}

// Run walks the repository and counts tokens by file and directory.
func Run(opts Options) (*Result, error) {
	if opts.Tokenizer == nil {
		return nil, fmt.Errorf("tokenizer is required")
	}

	root, err := filepath.Abs(opts.Root)
	if err != nil {
		return nil, fmt.Errorf("resolve root path: %w", err)
	}

	if opts.MaxFileBytes <= 0 {
		opts.MaxFileBytes = defaultMaxFileBytes
	}

	result := &Result{
		Repository:      root,
		Tokenizer:       opts.Tokenizer.Name(),
		TokenizerDetail: opts.Tokenizer.Description(),
		DirectoryTokens: map[string]int{".": 0},
	}

	err = filepath.WalkDir(root, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			if errors.Is(walkErr, os.ErrNotExist) || errors.Is(walkErr, os.ErrPermission) {
				return nil
			}
			return walkErr
		}
		if path == root {
			return nil
		}

		isDir := d.IsDir()
		if opts.IgnoreSpec != nil && opts.IgnoreSpec.MatchPath(path, isDir) {
			if isDir {
				result.IgnoredFiles += countFilesUnderDir(path)
				return fs.SkipDir
			}
			result.IgnoredFiles++
			return nil
		}

		if isDir {
			return nil
		}

		info, err := d.Info()
		if err != nil {
			if errors.Is(err, os.ErrNotExist) || errors.Is(err, os.ErrPermission) {
				return nil
			}
			return err
		}

		if info.Size() > opts.MaxFileBytes {
			result.IgnoredFiles++
			return nil
		}

		data, err := os.ReadFile(path)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) || errors.Is(err, os.ErrPermission) {
				return nil
			}
			return err
		}

		if isLikelyBinary(data) {
			result.IgnoredFiles++
			return nil
		}

		tokens := opts.Tokenizer.Count(string(data))
		result.TotalTokens += tokens
		result.TotalFiles++
		result.TotalLines += countLines(data)

		relPath, err := filepath.Rel(root, path)
		if err == nil {
			addTokensToDirs(result.DirectoryTokens, relPath, tokens)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("walk repository: %w", err)
	}

	result.DirectoryTokens["."] = result.TotalTokens
	return result, nil
}

func addTokensToDirs(dirTotals map[string]int, relPath string, tokens int) {
	relPath = filepath.Clean(relPath)
	relDir := filepath.Dir(relPath)
	if relDir == "" {
		relDir = "."
	}

	for {
		if relDir == "" {
			relDir = "."
		}
		dirTotals[relDir] += tokens
		if relDir == "." {
			break
		}
		relDir = filepath.Dir(relDir)
	}
}

func countFilesUnderDir(dir string) int {
	total := 0
	_ = filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			if errors.Is(err, os.ErrNotExist) || errors.Is(err, os.ErrPermission) {
				return nil
			}
			return nil
		}
		if d.IsDir() {
			return nil
		}
		total++
		return nil
	})
	return total
}

func isLikelyBinary(data []byte) bool {
	if len(data) == 0 {
		return false
	}
	if bytes.IndexByte(data, 0) >= 0 {
		return true
	}

	nonText := 0
	limit := len(data)
	if limit > 4096 {
		limit = 4096
	}
	for _, b := range data[:limit] {
		if (b < 9) || (b > 13 && b < 32) {
			nonText++
		}
	}
	return float64(nonText)/float64(limit) > 0.30
}

func countLines(data []byte) int {
	if len(data) == 0 {
		return 0
	}
	return bytes.Count(data, []byte{'\n'}) + 1
}
