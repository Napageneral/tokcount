package tokenizer

import (
	"fmt"
	"strings"
)

// Tokenizer counts tokens for plain text.
type Tokenizer interface {
	Count(text string) int
	Name() string
	Description() string
}

// New creates a tokenizer by name.
func New(name string) (Tokenizer, error) {
	switch normalize(name) {
	case "", "estimate", "google", "gemini":
		return NewEstimate(3.5), nil
	case "anthropic", "claude":
		return NewTiktoken("anthropic", "cl100k_base", "cl100k_base (Claude approximation)")
	case "openai":
		return NewTiktoken("openai", "cl100k_base", "cl100k_base (GPT-4)")
	case "openai-o200k":
		return NewTiktoken("openai-o200k", "o200k_base", "o200k_base (GPT-4o/o1)")
	default:
		return nil, fmt.Errorf("unsupported tokenizer: %s", name)
	}
}

func normalize(name string) string {
	return strings.ToLower(strings.TrimSpace(name))
}
