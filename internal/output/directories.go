package output

import (
	"path/filepath"
	"sort"

	"github.com/Napageneral/tokcount/internal/count"
)

// DirectoryStat is a normalized directory-level token rollup.
type DirectoryStat struct {
	Path       string  `json:"path"`
	Tokens     int     `json:"tokens"`
	Percentage float64 `json:"percentage"`
}

// AllDirectoryStats returns all non-root directories sorted by tokens desc.
func AllDirectoryStats(result *count.Result) []DirectoryStat {
	if result == nil || result.TotalTokens <= 0 {
		return nil
	}
	stats := make([]DirectoryStat, 0, len(result.DirectoryTokens))
	for path, tokens := range result.DirectoryTokens {
		if path == "." || tokens <= 0 {
			continue
		}
		pct := (float64(tokens) / float64(result.TotalTokens)) * 100
		stats = append(stats, DirectoryStat{
			Path:       normalizeDirectoryPath(path),
			Tokens:     tokens,
			Percentage: pct,
		})
	}
	sort.Slice(stats, func(i, j int) bool {
		if stats[i].Tokens == stats[j].Tokens {
			return stats[i].Path < stats[j].Path
		}
		return stats[i].Tokens > stats[j].Tokens
	})
	return stats
}

// TopDirectoryStats returns top N directory rows and remaining count.
func TopDirectoryStats(result *count.Result, limit int) ([]DirectoryStat, int) {
	all := AllDirectoryStats(result)
	if limit <= 0 || limit >= len(all) {
		return all, 0
	}
	return all[:limit], len(all) - limit
}

func normalizeDirectoryPath(path string) string {
	path = filepath.ToSlash(path)
	if path == "." || path == "" {
		return "./"
	}
	return path + "/"
}
