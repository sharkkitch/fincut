// Package reorder provides line reordering by field-based or reverse ordering.
package reorder

import (
	"fmt"
	"sort"
	"strings"
)

// Options configures the Reorderer.
type Options struct {
	// Reverse reverses the input order.
	Reverse bool
	// Field is the 1-based field index to sort by (0 means whole line).
	Field int
	// Delimiter separates fields when Field > 0.
	Delimiter string
	// Stable uses a stable sort algorithm.
	Stable bool
}

// Reorderer reorders lines according to configured options.
type Reorderer struct {
	opts Options
}

// New creates a Reorderer from opts, returning an error if opts are invalid.
func New(opts Options) (*Reorderer, error) {
	if err := validateOptions(opts); err != nil {
		return nil, err
	}
	return &Reorderer{opts: opts}, nil
}

// Apply reorders lines and returns the result.
func (r *Reorderer) Apply(lines []string) []string {
	out := make([]string, len(lines))
	copy(out, lines)

	if r.opts.Reverse && r.opts.Field == 0 {
		for i, j := 0, len(out)-1; i < j; i, j = i+1, j-1 {
			out[i], out[j] = out[j], out[i]
		}
		return out
	}

	key := func(line string) string {
		if r.opts.Field == 0 || r.opts.Delimiter == "" {
			return line
		}
		parts := strings.Split(line, r.opts.Delimiter)
		idx := r.opts.Field - 1
		if idx >= len(parts) {
			return ""
		}
		return parts[idx]
	}

	less := func(i, j int) bool {
		ki, kj := key(out[i]), key(out[j])
		if r.opts.Reverse {
			return ki > kj
		}
		return ki < kj
	}

	if r.opts.Stable {
		sort.SliceStable(out, less)
	} else {
		sort.Slice(out, less)
	}
	return out
}

// FormatSummary returns a human-readable summary of reorder results.
func FormatSummary(in, out []string) string {
	return fmt.Sprintf("reordered %d lines", len(out))
}
