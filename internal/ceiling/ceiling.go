// Package ceiling provides a line limiter that drops lines once a per-pattern
// match ceiling is reached.
package ceiling

import (
	"fmt"
	"regexp"
)

// Options configures the Ceiling limiter.
type Options struct {
	// Pattern is the regex to match against each line.
	Pattern string
	// Max is the maximum number of matching lines to allow through.
	Max int
}

// Ceiling drops matching lines once the match count hits Max.
type Ceiling struct {
	re    *regexp.Regexp
	max   int
	count int
}

// New constructs a Ceiling from opts.
func New(opts Options) (*Ceiling, error) {
	if opts.Pattern == "" {
		return nil, fmt.Errorf("ceiling: pattern must not be empty")
	}
	if opts.Max <= 0 {
		return nil, fmt.Errorf("ceiling: max must be greater than zero, got %d", opts.Max)
	}
	re, err := regexp.Compile(opts.Pattern)
	if err != nil {
		return nil, fmt.Errorf("ceiling: invalid pattern: %w", err)
	}
	return &Ceiling{re: re, max: opts.Max}, nil
}

// Apply filters lines, passing through non-matching lines unconditionally and
// matching lines only until Max is reached.
func (c *Ceiling) Apply(lines []string) []string {
	out := make([]string, 0, len(lines))
	for _, l := range lines {
		if c.re.MatchString(l) {
			if c.count >= c.max {
				continue
			}
			c.count++
		}
		out = append(out, l)
	}
	return out
}

// FormatSummary returns a human-readable summary of the ceiling state.
func (c *Ceiling) FormatSummary() string {
	return fmt.Sprintf("ceiling: %d/%d matching lines passed", c.count, c.max)
}
