package sort

import (
	"fmt"
	"sort"
	"strings"
)

// Sorter sorts lines lexicographically or by a specific field delimiter.
type Sorter struct {
	reverse   bool
	unique    bool
	delimiter string
	field     int
}

// Options configures a Sorter.
type Options struct {
	Reverse   bool
	Unique    bool
	Delimiter string
	Field     int // 1-based; 0 means sort whole line
}

// New creates a Sorter from Options.
func New(opts Options) (*Sorter, error) {
	if opts.Field < 0 {
		return nil, fmt.Errorf("sort: field index must be >= 0, got %d", opts.Field)
	}
	if opts.Field > 0 && opts.Delimiter == "" {
		return nil, fmt.Errorf("sort: delimiter required when field > 0")
	}
	return &Sorter{
		reverse:   opts.Reverse,
		unique:    opts.Unique,
		delimiter: opts.Delimiter,
		field:     opts.Field,
	}, nil
}

// Apply sorts the provided lines and returns the result.
func (s *Sorter) Apply(lines []string) []string {
	if len(lines) == 0 {
		return lines
	}

	out := make([]string, len(lines))
	copy(out, lines)

	sort.SliceStable(out, func(i, j int) bool {
		a := s.key(out[i])
		b := s.key(out[j])
		if s.reverse {
			return a > b
		}
		return a < b
	})

	if s.unique {
		out = deduplicate(out, s.key)
	}

	return out
}

func (s *Sorter) key(line string) string {
	if s.field == 0 || s.delimiter == "" {
		return line
	}
	parts := strings.Split(line, s.delimiter)
	idx := s.field - 1
	if idx >= len(parts) {
		return ""
	}
	return parts[idx]
}

func deduplicate(lines []string, key func(string) string) []string {
	seen := make(map[string]struct{}, len(lines))
	out := lines[:0]
	for _, l := range lines {
		k := key(l)
		if _, ok := seen[k]; !ok {
			seen[k] = struct{}{}
			out = append(out, l)
		}
	}
	return out
}
