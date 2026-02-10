package tokenizer

import (
	"fmt"

	tiktoken "github.com/pkoukk/tiktoken-go"
)

// TiktokenTokenizer uses a concrete tiktoken encoding.
type TiktokenTokenizer struct {
	name        string
	description string
	encoder     *tiktoken.Tiktoken
}

// NewTiktoken creates a tokenizer backed by a tiktoken encoding.
func NewTiktoken(name string, encoding string, description string) (*TiktokenTokenizer, error) {
	encoder, err := tiktoken.GetEncoding(encoding)
	if err != nil {
		return nil, fmt.Errorf("init tokenizer %s (%s): %w", name, encoding, err)
	}
	return &TiktokenTokenizer{
		name:        name,
		description: description,
		encoder:     encoder,
	}, nil
}

func (t *TiktokenTokenizer) Count(text string) int {
	if text == "" {
		return 0
	}
	tokens := t.encoder.Encode(text, nil, nil)
	return len(tokens)
}

func (t *TiktokenTokenizer) Name() string {
	return t.name
}

func (t *TiktokenTokenizer) Description() string {
	return t.description
}
