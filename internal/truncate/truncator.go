// Package truncate provides utilities for truncating long log lines
// to a configurable maximum byte or rune length, with optional ellipsis.
package truncate

import (
	"fmt"
	"unicode/utf8"
)

// Options configures the Truncator.
type Options struct {
	// MaxRunes is the maximum number of runes per line (0 = unlimited).
	MaxRunes int
	// Ellipsis is appended when a line is truncated. Defaults to "...".
	Ellipsis string
	// StripANSI removes ANSI escape sequences before measuring length.
	StripANSI bool
}

// Truncator truncates lines that exceed a configured rune limit.
type Truncator struct {
	opts Options
}

// NewTruncator creates a Truncator from the given Options.
// Returns an error if MaxRunes is negative.
func NewTruncator(opts Options) (*Truncator, error) {
	if opts.MaxRunes < 0 {
		return nil, fmt.Errorf("truncate: MaxRunes must be >= 0, got %d", opts.MaxRunes)
	}
	if opts.Ellipsis == "" {
		opts.Ellipsis = "..."
	}
	return &Truncator{opts: opts}, nil
}

// Apply truncates each line in lines according to the configured options.
// Lines that are within the limit are returned unchanged.
func (t *Truncator) Apply(lines []string) []string {
	if t.opts.MaxRunes == 0 {
		out := make([]string, len(lines))
		copy(out, lines)
		return out
	}

	ellipsisRunes := utf8.RuneCountInString(t.opts.Ellipsis)
	cutAt := t.opts.MaxRunes - ellipsisRunes
	if cutAt < 0 {
		cutAt = 0
	}

	out := make([]string, 0, len(lines))
	for _, line := range lines {
		count := utf8.RuneCountInString(line)
		if count <= t.opts.MaxRunes {
			out = append(out, line)
			continue
		}
		// Trim to cutAt runes then append ellipsis.
		i := 0
		byte_idx := 0
		for _, r := range line {
			if i >= cutAt {
				break
			}
			byte_idx += utf8.RuneLen(r)
			i++
		}
		out = append(out, line[:byte_idx]+t.opts.Ellipsis)
	}
	return out
}
