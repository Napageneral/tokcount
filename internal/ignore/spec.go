package ignore

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	gitignore "github.com/sabhiram/go-gitignore"
)

// Spec encapsulates ignore matching rooted at a repository path.
type Spec struct {
	root     string
	matcher  *gitignore.GitIgnore
	patterns []string
}

// LoadSpec compiles default patterns + .cartographerignore + .gitignore
// and optionally a custom ignore file path.
func LoadSpec(scopeRoot string, customIgnoreFile string) (*Spec, error) {
	root, err := filepath.Abs(scopeRoot)
	if err != nil {
		return nil, fmt.Errorf("resolve root path: %w", err)
	}

	patterns := DefaultPatterns()

	if p, err := readIgnoreFileOptional(filepath.Join(root, ".cartographerignore")); err != nil {
		return nil, err
	} else {
		patterns = append(patterns, p...)
	}

	if p, err := readIgnoreFileOptional(filepath.Join(root, ".gitignore")); err != nil {
		return nil, err
	} else {
		patterns = append(patterns, p...)
	}

	if strings.TrimSpace(customIgnoreFile) != "" {
		customPath := customIgnoreFile
		if !filepath.IsAbs(customPath) {
			customPath = filepath.Join(root, customPath)
		}
		p, err := readIgnoreFileRequired(customPath)
		if err != nil {
			return nil, fmt.Errorf("load custom ignore file: %w", err)
		}
		patterns = append(patterns, p...)
	}

	patterns = dedupePatterns(patterns)

	var matcher *gitignore.GitIgnore
	if len(patterns) > 0 {
		matcher = gitignore.CompileIgnoreLines(patterns...)
	}

	return &Spec{
		root:     root,
		matcher:  matcher,
		patterns: patterns,
	}, nil
}

func (s *Spec) Root() string {
	if s == nil {
		return ""
	}
	return s.root
}

func (s *Spec) Patterns() []string {
	if s == nil {
		return nil
	}
	out := make([]string, len(s.patterns))
	copy(out, s.patterns)
	return out
}

// MatchPath reports whether an absolute path should be ignored.
func (s *Spec) MatchPath(absPath string, isDir bool) bool {
	if s == nil || s.matcher == nil {
		return false
	}
	rel, err := filepath.Rel(s.root, absPath)
	if err != nil {
		rel = absPath
	}
	rel = filepath.Clean(rel)
	if rel == "." || rel == "" {
		return false
	}

	rel = filepath.ToSlash(rel)
	if isDir && !strings.HasSuffix(rel, "/") {
		rel += "/"
	}
	return s.matcher.MatchesPath(rel)
}

func readIgnoreFileOptional(path string) ([]string, error) {
	_, err := os.Stat(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}
		return nil, err
	}
	return readIgnoreFile(path)
}

func readIgnoreFileRequired(path string) ([]string, error) {
	return readIgnoreFile(path)
}

func readIgnoreFile(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	out := make([]string, 0)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "#") {
			continue
		}
		out = append(out, line)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return out, nil
}

func dedupePatterns(patterns []string) []string {
	seen := make(map[string]bool, len(patterns))
	out := make([]string, 0, len(patterns))
	for _, p := range patterns {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		if seen[p] {
			continue
		}
		seen[p] = true
		out = append(out, p)
	}
	return out
}
