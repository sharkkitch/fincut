// Package squeeze provides a processor that collapses consecutive repeated
// lines into a single occurrence, similar to the Unix `uniq` command but with
// optional repeat-count annotation.
package squeeze

import "fmt"

// Options configures the Squeezer.
type Options struct {
	// Annotate, when true, prepends the repeat count to each emitted line.
	Annotate bool
	// Separator is placed between the count annotation and the line content.
	// Defaults to "\t" when Annotate is true and Separator is empty.
	Separator string
}

// Squeezer collapses consecutive duplicate lines.
type Squeezer struct {
	opts Options
}

// New creates a new Squeezer with the provided options.
func New(opts Options) (*Squeezer, error) {
	if opts.Annotate && opts.Separator == "" {
		opts.Separator = "\t"
	}
	return &Squeezer{opts: opts}, nil
}

// Apply processes lines and returns a new slice with consecutive duplicates
// removed. If Annotate is enabled each output line is prefixed with its run
// length.
func (s *Squeezer) Apply(lines []string) []string {
	if len(lines) == 0 {
		return []string{}
	}

	out := make([]string, 0, len(lines))
	current := lines[0]
	count := 1

	flush := func() {
		if s.opts.Annotate {
			out = append(out, fmt.Sprintf("%d%s%s", count, s.opts.Separator, current))
		} else {
			out = append(out, current)
		}
	}

	for _, line := range lines[1:] {
		if line == current {
			count++
			continue
		}
		flush()
		current = line
		count = 1
	}
	flush()

	return out
}

// FormatSummary returns a human-readable summary of how many lines were
// removed by squeezing.
func FormatSummary(input, output int) string {
	return fmt.Sprintf("squeeze: %d → %d lines (%d removed)", input, output, input-output)
}
