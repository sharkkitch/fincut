// Package unique provides line-level uniqueness filtering by a configurable key field.
package unique

import (
	"fmt"
	"strings"
)

// Options configures the Uniquer.
type Options struct {
	// Field is the 1-based field index to use as the uniqueness key.
	// If 0, the entire line is used.
	Field int
	// Delimiter separates fields when Field > 0.
	Delimiter string
	// CaseInsensitive treats keys as case-insensitive.
	CaseInsensitive bool
}

// Uniquer filters lines so that only the first occurrence of each key is kept.
type Uniquer struct {
	opts Options
}

// New creates a new Uniquer with the given options.
func New(opts Options) (*Uniquer, error) {
	if err := validateOptions(opts); err != nil {
		return nil, err
	}
	return &Uniquer{opts: opts}, nil
}

func validateOptions(opts Options) error {
	if opts.Field < 0 {
		return fmt.Errorf("unique: field index must be >= 0, got %d", opts.Field)
	}
	if opts.Field > 0 && opts.Delimiter == "" {
		return fmt.Errorf("unique: delimiter required when field > 0")
	}
	return nil
}

// Apply filters lines, returning only the first occurrence of each unique key.
func (u *Uniquer) Apply(lines []string) []string {
	seen := make(map[string]struct{}, len(lines))
	out := make([]string, 0, len(lines))
	for _, line := range lines {
		key := u.extractKey(line)
		if u.opts.CaseInsensitive {
			key = strings.ToLower(key)
		}
		if _, exists := seen[key]; exists {
			continue
		}
		seen[key] = struct{}{}
		out = append(out, line)
	}
	return out
}

func (u *Uniquer) extractKey(line string) string {
	if u.opts.Field == 0 {
		return line
	}
	parts := strings.Split(line, u.opts.Delimiter)
	idx := u.opts.Field - 1
	if idx >= len(parts) {
		return line
	}
	return parts[idx]
}
