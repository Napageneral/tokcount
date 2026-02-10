package output

import (
	"encoding/json"
	"math"

	"github.com/Napageneral/tokcount/internal/count"
)

type jsonPayload struct {
	Repository      string          `json:"repository"`
	Tokenizer       string          `json:"tokenizer"`
	TotalTokens     int             `json:"total_tokens"`
	TotalFiles      int             `json:"total_files"`
	IgnoredFiles    int             `json:"ignored_files"`
	Directories     []DirectoryStat `json:"directories"`
	PricingEstimate PricingEstimate `json:"pricing_estimate"`
}

// RenderJSON marshals machine-readable token count output.
func RenderJSON(result *count.Result) ([]byte, error) {
	all := AllDirectoryStats(result)
	for i := range all {
		all[i].Percentage = math.Round(all[i].Percentage*10) / 10
	}

	payload := jsonPayload{
		Repository:      result.Repository,
		Tokenizer:       result.Tokenizer,
		TotalTokens:     result.TotalTokens,
		TotalFiles:      result.TotalFiles,
		IgnoredFiles:    result.IgnoredFiles,
		Directories:     all,
		PricingEstimate: EstimatePricing(result.TotalTokens),
	}
	payload.PricingEstimate.TokensMillions = math.Round(payload.PricingEstimate.TokensMillions*100) / 100

	return json.MarshalIndent(payload, "", "  ")
}
