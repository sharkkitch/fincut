// Package wrap provides line-wrapping for long log lines.
package wrap

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

// Options configures the Wrapper.
type Options struct {
	// Width is the maximum rune width per line. Must be > 0.
	Width int
	// Indent is prepended to continuation lines.
	Indent string
	// HardBreak forces a break at exactly Width even inside a word.
	// When false the break occurs at the last space before Width.
	HardBreak bool
}

// Wrapper wraps lines that exceed a configured width.
type Wrapper struct {
	opts Options
}

// New creates a Wrapper from opts. Returns an error if Width <= 0.
func New(opts Options) (*Wrapper, error) {
	if opts.Width <= 0 {
		return nil, fmt.Errorf("wrap: Width must be greater than zero, got %d", opts.Width)
	}
	return &Wrapper{opts: opts}, nil
}

// Apply wraps each line in lines and returns the expanded slice.
func (w *Wrapper) Apply(lines []string) []string {
	out := make([]string, 0, len(lines))
	for _, line := range lines {
		out = append(out, w.wrapLine(line)...)
	}
	return out
}

// wrapLine splits a single line into wrapped segments.
func (w *Wrapper) wrapLine(line string) []string {
	if utf8.RuneCountInString(line) <= w.opts.Width {
		return []string{line}
	}

	var result []string
	indent := ""
	remaining := line

	for utf8.RuneCountInString(remaining) > w.opts.Width {
		avail := w.opts.Width - utf8.RuneCountInString(indent)
		if avail <= 0 {
			// indent wider than width — emit remainder as-is to avoid infinite loop
			break
		}

		runes := []rune(remaining)
		breakAt := avail

		if !w.opts.HardBreak {
			// find last space at or before avail
			for i := avail - 1; i > 0; i-- {
				if runes[i] == ' ' {
					breakAt = i
					break
				}
			}
		}

		result = append(result, indent+string(runes[:breakAt]))
		remaining = strings.TrimLeft(string(runes[breakAt:]), " ")
		indent = w.opts.Indent
	}

	if remaining != "" {
		result = append(result, indent+remaining)
	}
	return result
}

// FormatSummary returns a one-line summary of wrap statistics.
func FormatSummary(in, out int) string {
	return fmt.Sprintf("wrap: %d lines in → %d lines out (%+d)", in, out, out-in)
}
