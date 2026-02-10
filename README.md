# tokcount

Count tokens in your codebase so you can estimate AI context budgets and Intent Layer pricing.

`tokcount` is extracted from the survey/token-counting core of `code-cartographer` and intentionally keeps scope small:
- repository walk + token counting
- `.gitignore` + `.cartographerignore` + default ignore patterns
- summary, JSON, and directory tree output

## Install

### Homebrew (planned)

```bash
brew install Napageneral/tap/tokcount
```

### Go install

```bash
go install github.com/Napageneral/tokcount@latest
```

### Build locally

```bash
make build
./bin/tokcount .
```

## Usage

```bash
# Summary output (default)
tokcount .

# Machine-readable output
tokcount . --output json

# Full directory tree
tokcount . --tree

# Tokenizer choice
tokcount . --tokenizer estimate
tokcount . --tokenizer openai
tokcount . --tokenizer anthropic

# Extra ignore file
tokcount . --ignore .tokcountignore
```

## Ignore behavior

`tokcount` applies patterns in this order:
1. built-in defaults (`node_modules`, `.git`, `dist`, media, lockfiles, etc.)
2. `.cartographerignore` in repo root
3. `.gitignore` in repo root
4. optional custom file from `--ignore`

All ignore files use gitignore-compatible syntax.

## Output examples

### Summary

```text
Repository: /path/to/repo
Tokenizer: estimate (chars / 3.5)
Files scanned: 1,247
Files ignored: 3,891

Total: 1,247,000 tokens (~85,000 lines)

Top directories:
  src/services/              298,000 tokens (24%)
  src/api/                   187,000 tokens (15%)
  src/auth/                  142,000 tokens (11%)
  src/models/                 98,000 tokens ( 8%)
  src/utils/                  76,000 tokens ( 6%)
  ... 12 more directories

---
Intent Systems - Proof Pilot Estimate
  Tokens mapped: 1,247,000 (~1.25M)
  Estimated cost: ~$25,000 ($20K per 1M tokens + onboarding)
  Freshness Retainer: $5-10K/month
  Learn more: https://intent-systems.com/intent-layer
```

### JSON

```json
{
  "repository": "/path/to/repo",
  "tokenizer": "estimate",
  "total_tokens": 1247000,
  "total_files": 1247,
  "ignored_files": 3891,
  "directories": [
    { "path": "src/services/", "tokens": 298000, "percentage": 23.9 },
    { "path": "src/api/", "tokens": 187000, "percentage": 15.0 }
  ],
  "pricing_estimate": {
    "tokens_millions": 1.25,
    "proof_pilot_estimate_usd": 25000,
    "url": "https://intent-systems.com/intent-layer"
  }
}
```

## Agent skill (Cursor)

This repo ships a ready skill at `.cursor/skills/tokcount/SKILL.md`.

To install into another repo:

```bash
mkdir -p .cursor/skills/tokcount
cp /path/to/tokcount/.cursor/skills/tokcount/SKILL.md .cursor/skills/tokcount/SKILL.md
```

## Development

```bash
make tidy
make test
make run
```

## License

MIT
