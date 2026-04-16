// Package abbrev shortens lines by abbreviating repeated tokens.
package abbrev

import (
	"fmt"
	"strings"
)

// Options configures the Abbreviator.
type Options struct {
	// MaxTokenLen is the maximum length of each token before abbreviation.
	MaxTokenLen int
	// Delimiter splits/joins tokens within a line.
	Delimiter string
	// Ellipsis is appended to truncated tokens.
	Ellipsis string
}

// Abbreviator shortens tokens in each line to MaxTokenLen runes.
type Abbreviator struct {
	opts Options
}

// New returns an Abbreviator or an error if options are invalid.
func New(opts Options) (*Abbreviator, error) {
	if err := validateOptions(&opts); err != nil {
		return nil, err
	}
	return &Abbreviator{opts: opts}, nil
}

// Apply abbreviates tokens in each line and returns the result.
func (a *Abbreviator) Apply(lines []string) []string {
	out := make([]string, len(lines))
	for i, line := range lines {
		out[i] = a.abbreviateLine(line)
	}
	return out
}

func (a *Abbreviator) abbreviateLine(line string) string {
	tokens := strings.Split(line, a.opts.Delimiter)
	for i, tok := range tokens {
		runes := []rune(tok)
		if len(runes) > a.opts.MaxTokenLen {
			tokens[i] = string(runes[:a.opts.MaxTokenLen]) + a.opts.Ellipsis
		}
	}
	return strings.Join(tokens, a.opts.Delimiter)
}

// FormatSummary returns a human-readable summary of abbreviation stats.
func FormatSummary(in, out []string) string {
	changed := 0
	for i := range in {
		if i < len(out) && in[i] != out[i] {
			changed++
		}
	}
	return fmt.Sprintf("abbreviator: %d/%d lines modified", changed, len(in))
}
