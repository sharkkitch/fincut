// Package fence extracts lines between paired delimiter patterns,
// optionally including or excluding the delimiter lines themselves.
package fence

import (
	"fmt"
	"regexp"
)

// Options configures the Fencer.
type Options struct {
	OpenPattern  string
	ClosePattern string
	IncludeDelim bool
}

// Fencer extracts fenced blocks from a slice of lines.
type Fencer struct {
	open         *regexp.Regexp
	close        *regexp.Regexp
	includeDelim bool
}

// New constructs a Fencer from the given Options.
func New(opts Options) (*Fencer, error) {
	if err := validateOptions(opts); err != nil {
		return nil, err
	}
	op, err := regexp.Compile(opts.OpenPattern)
	if err != nil {
		return nil, fmt.Errorf("fence: invalid open pattern: %w", err)
	}
	cl, err := regexp.Compile(opts.ClosePattern)
	if err != nil {
		return nil, fmt.Errorf("fence: invalid close pattern: %w", err)
	}
	return &Fencer{open: op, close: cl, includeDelim: opts.IncludeDelim}, nil
}

// Apply returns only the lines that fall within fenced blocks.
// Multiple non-overlapping blocks are all collected.
func (f *Fencer) Apply(lines []string) []string {
	var out []string
	inside := false
	for _, line := range lines {
		if !inside && f.open.MatchString(line) {
			inside = true
			if f.includeDelim {
				out = append(out, line)
			}
			continue
		}
		if inside && f.close.MatchString(line) {
			inside = false
			if f.includeDelim {
				out = append(out, line)
			}
			continue
		}
		if inside {
			out = append(out, line)
		}
	}
	return out
}

// FormatSummary returns a human-readable summary of extraction results.
func FormatSummary(total, extracted int) string {
	return fmt.Sprintf("fence: %d/%d lines extracted from fenced blocks", extracted, total)
}
