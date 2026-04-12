package join

import (
	"fmt"
	"strings"
)

// Joiner merges consecutive lines into single lines using a configurable separator.
type Joiner struct {
	opts Options
}

// New creates a new Joiner with the given options.
func New(opts Options) (*Joiner, error) {
	if err := validateOptions(&opts); err != nil {
		return nil, err
	}
	return &Joiner{opts: opts}, nil
}

// Apply joins input lines into groups of GroupSize, separated by Separator.
// If GroupSize is 0, all lines are joined into a single line.
func (j *Joiner) Apply(lines []string) []string {
	if len(lines) == 0 {
		return []string{}
	}

	if j.opts.GroupSize == 0 {
		return []string{strings.Join(lines, j.opts.Separator)}
	}

	result := make([]string, 0, (len(lines)+j.opts.GroupSize-1)/j.opts.GroupSize)
	for i := 0; i < len(lines); i += j.opts.GroupSize {
		end := i + j.opts.GroupSize
		if end > len(lines) {
			end = len(lines)
		}
		result = append(result, strings.Join(lines[i:end], j.opts.Separator))
	}
	return result
}

// FormatSummary returns a human-readable summary of a join operation.
func FormatSummary(inputCount, outputCount int) string {
	return fmt.Sprintf("joined %d lines into %d records", inputCount, outputCount)
}
