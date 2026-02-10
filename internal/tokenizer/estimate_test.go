package tokenizer

import "testing"

func TestEstimateTokenizer_Count(t *testing.T) {
	tok := NewEstimate(3.5)

	got := tok.Count("1234567")
	want := 2
	if got != want {
		t.Fatalf("expected %d tokens, got %d", want, got)
	}
}

func TestFactory_DefaultEstimate(t *testing.T) {
	tok, err := New("")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if tok.Name() != "estimate" {
		t.Fatalf("expected estimate tokenizer, got %s", tok.Name())
	}
}
