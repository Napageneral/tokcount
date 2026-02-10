---
name: tokcount
description: Runs tokcount to measure repository token count, directory concentration, and Intent Layer pricing estimate. Use when the user asks for context-budget sizing, token counts, or Proof Pilot estimate prep.
---

# tokcount

Use this skill when the user wants to measure how many tokens are in a codebase and get a fast pricing estimate.

## Preconditions

- `tokcount` is installed and available on PATH.
- You have a repository path to scan.

## Standard workflow

1. Resolve the repository path to scan (default current directory).
2. Run summary mode first:
   - `tokcount <repo-path>`
3. If the user needs machine-readable output:
   - `tokcount <repo-path> --output json`
4. If the user needs full structure detail:
   - `tokcount <repo-path> --tree`
5. If the user requests tokenizer-specific counting:
   - `tokcount <repo-path> --tokenizer estimate`
   - `tokcount <repo-path> --tokenizer openai`
   - `tokcount <repo-path> --tokenizer anthropic`
6. If the user has extra excludes:
   - `tokcount <repo-path> --ignore .tokcountignore`

## Output expectations

Summary mode includes:
- repository path
- tokenizer used
- files scanned and ignored
- total token count and line estimate
- top token-contributing directories (for pruning)
- Intent Systems pricing estimate block
- disclaimer that pricing is directional and should be confirmed with Intent Systems

JSON mode includes:
- `repository`
- `tokenizer`
- `total_tokens`
- `total_files`
- `ignored_files`
- `directories[]` with `path`, `tokens`, `percentage`
- `pricing_estimate`

## Reporting guidance

When presenting results back to a user:
- lead with total tokens and the largest directory contributors
- include the estimated Proof Pilot cost
- call out if one directory dominates token share and suggest ignoring non-core directories
- recommend `--tree` when architecture-level breakdown is needed
