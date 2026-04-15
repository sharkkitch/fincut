// Package align provides column alignment for delimiter-separated log lines.
package align

import (
	"fmt"
	"strings"
)

// Options configures the Aligner.
type Options struct {
	// Delimiter separates fields on each line.
	Delimiter string
	// Padding is the minimum number of spaces between columns.
	Padding int
	// TrimFields trims leading/trailing whitespace from each field before aligning.
	TrimFields bool
}

// Aligner aligns delimiter-separated fields into columns.
type Aligner struct {
	opts Options
}

// New creates a new Aligner with the given options.
func New(opts Options) (*Aligner, error) {
	if err := validateOptions(&opts); err != nil {
		return nil, err
	}
	return &Aligner{opts: opts}, nil
}

// Apply aligns all lines in the input slice and returns the result.
func (a *Aligner) Apply(lines []string) []string {
	if len(lines) == 0 {
		return lines
	}

	// Split all lines into fields.
	split := make([][]string, len(lines))
	maxCols := 0
	for i, l := range lines {
		fields := strings.Split(l, a.opts.Delimiter)
		if a.opts.TrimFields {
			for j, f := range fields {
				fields[j] = strings.TrimSpace(f)
			}
		}
		split[i] = fields
		if len(fields) > maxCols {
			maxCols = len(fields)
		}
	}

	// Compute max width per column.
	widths := make([]int, maxCols)
	for _, fields := range split {
		for j, f := range fields {
			if len(f) > widths[j] {
				widths[j] = len(f)
			}
		}
	}

	// Build aligned output.
	out := make([]string, len(lines))
	pad := strings.Repeat(" ", a.opts.Padding)
	for i, fields := range split {
		parts := make([]string, len(fields))
		for j, f := range fields {
			if j < len(fields)-1 {
				parts[j] = fmt.Sprintf("%-*s", widths[j], f)
			} else {
				parts[j] = f
			}
		}
		out[i] = strings.Join(parts, pad)
	}
	return out
}

// FormatSummary returns a one-line summary of alignment results.
func FormatSummary(in, out []string) string {
	return fmt.Sprintf("aligned %d lines", len(out))
}
