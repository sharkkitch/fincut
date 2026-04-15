// Package fold provides line-folding by merging continuation lines
// into a single logical line based on a configurable pattern.
package fold

import (
	"fmt"
	"regexp"
)

// Options configures the Folder.
type Options struct {
	// ContinuationPattern is matched against a line to decide whether it
	// is a continuation of the previous line. Required.
	ContinuationPattern string

	// Separator is placed between merged lines. Defaults to " ".
	Separator string
}

// Folder merges continuation lines into their predecessor.
type Folder struct {
	re        *regexp.Regexp
	separator string
}

// New creates a Folder from opts.
func New(opts Options) (*Folder, error) {
	if err := validateOptions(&opts); err != nil {
		return nil, err
	}
	re, err := regexp.Compile(opts.ContinuationPattern)
	if err != nil {
		return nil, fmt.Errorf("fold: invalid continuation pattern: %w", err)
	}
	return &Folder{re: re, separator: opts.Separator}, nil
}

// Apply folds continuation lines in lines and returns the result.
func (f *Folder) Apply(lines []string) []string {
	if len(lines) == 0 {
		return []string{}
	}
	out := make([]string, 0, len(lines))
	current := lines[0]
	for _, line := range lines[1:] {
		if f.re.MatchString(line) {
			current += f.separator + line
		} else {
			out = append(out, current)
			current = line
		}
	}
	out = append(out, current)
	return out
}

// FormatSummary returns a one-line description of the fold operation.
func FormatSummary(in, out int) string {
	return fmt.Sprintf("fold: %d lines -> %d lines (%d merged)", in, out, in-out)
}
