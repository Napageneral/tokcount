package tokenizer

import "math"

// EstimateTokenizer uses a characters-per-token heuristic.
type EstimateTokenizer struct {
	charsPerToken float64
}

// NewEstimate creates the chars/token estimator.
func NewEstimate(charsPerToken float64) *EstimateTokenizer {
	if charsPerToken <= 0 {
		charsPerToken = 3.5
	}
	return &EstimateTokenizer{charsPerToken: charsPerToken}
}

func (t *EstimateTokenizer) Count(text string) int {
	if text == "" {
		return 0
	}
	return int(math.Round(float64(len(text)) / t.charsPerToken))
}

func (t *EstimateTokenizer) Name() string {
	return "estimate"
}

func (t *EstimateTokenizer) Description() string {
	return "estimate (chars / 3.5)"
}
