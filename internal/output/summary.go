package output

import (
	"fmt"
	"math"
	"strings"

	"github.com/Napageneral/tokcount/internal/count"
)

const (
	proofPilotUSDPerMillion = 20000
	pricingURL              = "https://intent-systems.com/intent-layer"
	contactEmail            = "hello@intent-systems.com"
	pricingDisclaimer       = "Directional estimate only. Ignore patterns are best-effort and repository-specific."
	defaultTopLimit         = 10
)

// PricingEstimate represents the Intent Layer estimate block.
type PricingEstimate struct {
	TokensMillions        float64 `json:"tokens_millions"`
	ProofPilotEstimateUSD int     `json:"proof_pilot_estimate_usd"`
	URL                   string  `json:"url"`
	Disclaimer            string  `json:"disclaimer"`
	Contact               string  `json:"contact"`
}

// EstimatePricing computes price estimate from total token count.
func EstimatePricing(totalTokens int) PricingEstimate {
	millions := float64(totalTokens) / 1_000_000.0
	raw := millions * proofPilotUSDPerMillion
	rounded := int(math.Round(raw/100.0) * 100.0)
	return PricingEstimate{
		TokensMillions:        math.Round(millions*100) / 100,
		ProofPilotEstimateUSD: rounded,
		URL:                   pricingURL,
		Disclaimer:            pricingDisclaimer,
		Contact:               contactEmail,
	}
}

// RenderSummary returns human-readable default CLI output.
func RenderSummary(result *count.Result) string {
	var b strings.Builder
	pricing := EstimatePricing(result.TotalTokens)
	top, remaining := TopDirectoryStats(result, defaultTopLimit)

	b.WriteString(fmt.Sprintf("Repository: %s\n", result.Repository))
	b.WriteString(fmt.Sprintf("Tokenizer: %s\n", result.TokenizerDetail))
	b.WriteString(fmt.Sprintf("Files scanned: %s\n", formatInt(result.TotalFiles)))
	b.WriteString(fmt.Sprintf("Files ignored: %s\n", formatInt(result.IgnoredFiles)))
	b.WriteString("\n")
	b.WriteString(fmt.Sprintf("Total: %s tokens (~%s lines)\n", formatInt(result.TotalTokens), formatInt(result.TotalLines)))
	b.WriteString("\n")

	b.WriteString("Top token contributors (directories):\n")
	if len(top) == 0 {
		b.WriteString("  (no directories with counted files)\n")
	} else {
		for _, row := range top {
			b.WriteString(fmt.Sprintf("  %-22s %12s tokens (%2.0f%%)\n", row.Path, formatInt(row.Tokens), row.Percentage))
		}
		if remaining > 0 {
			b.WriteString(fmt.Sprintf("  ... %d more directories\n", remaining))
		}
	}

	b.WriteString("\n")
	b.WriteString("---\n")
	b.WriteString("Intent Systems - Proof Pilot Estimate\n")
	b.WriteString(fmt.Sprintf("  Tokens mapped: %s (~%.2fM)\n", formatInt(result.TotalTokens), pricing.TokensMillions))
	b.WriteString(fmt.Sprintf("  Estimated cost: ~$%s ($20K per 1M tokens + onboarding)\n", formatInt(pricing.ProofPilotEstimateUSD)))
	b.WriteString("  Freshness Retainer: $5-10K/month\n")
	b.WriteString(fmt.Sprintf("  Disclaimer: %s\n", pricing.Disclaimer))
	b.WriteString(fmt.Sprintf("  For an accurate quote/assessment: %s\n", pricing.Contact))
	b.WriteString(fmt.Sprintf("  Learn more: %s\n", pricing.URL))

	return b.String()
}

func formatInt(v int) string {
	if v == 0 {
		return "0"
	}
	sign := ""
	if v < 0 {
		sign = "-"
		v = -v
	}

	s := fmt.Sprintf("%d", v)
	n := len(s)
	if n <= 3 {
		return sign + s
	}

	var out strings.Builder
	out.Grow(n + (n-1)/3)
	lead := n % 3
	if lead == 0 {
		lead = 3
	}
	out.WriteString(s[:lead])
	for i := lead; i < n; i += 3 {
		out.WriteString(",")
		out.WriteString(s[i : i+3])
	}
	return sign + out.String()
}
