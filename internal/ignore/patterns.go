package ignore

// DefaultPatterns returns baseline ignore patterns.
func DefaultPatterns() []string {
	return []string{
		// VCS and tool state
		".git",
		".intent",
		".tldr",
		".DS_Store",

		// Dependency and build outputs
		"node_modules",
		"vendor",
		"dist",
		"build",
		".next",
		"__pycache__",
		".venv",
		"venv",
		"target",

		// Minified and bundled assets
		"*.min.js",
		"*.min.css",
		"*.bundle.js",

		// Language/tooling caches
		".mypy_cache",
		".pytest_cache",

		// Lock files
		"*.lock",
		"uv.lock",
		"pnpm-lock.yaml",
		"package-lock.json",
		"yarn.lock",
		"bun.lockb",
		"poetry.lock",
		"Pipfile.lock",
		"composer.lock",
		"Gemfile.lock",
		"Cargo.lock",

		// Compiled artifacts
		"*.so",
		"*.dylib",
		"*.a",
		"*.o",
		"*.pyc",
		"*.jar",
		"*.node",

		// Databases and blob-ish data
		"*.db",
		"*.db-wal",
		"*.db-shm",
		"*.db-journal",
		"*.sqlite",
		"*.sqlite3",
		"*.sqlite3-journal",
		"*.npz",
		"*.npy",
		"*.dat",
		"*.pkl",
		"*.sav",
		"*.csv",

		// Media and binary-rich assets
		"*.png",
		"*.jpg",
		"*.jpeg",
		"*.gif",
		"*.ico",
		"*.icns",
		"*.svg",
		"*.ttf",
		"*.webp",
		"*.wav",
		"*.mp3",
		"*.mp4",

		// Archives
		"*.gz",
		"*.bz2",
		"*.xz",
		"*.lzma",
		"*.zip",
		"*.tar",
		"*.tgz",
		"*.7z",
	}
}
