// Package linenum provides a processor that filters lines by line number ranges.
package linenum

import "fmt"

// Options configures the line number filter.
type Options struct {
	// Ranges is a list of [start, end] pairs (1-indexed, inclusive).
	// A range with end == 0 means "from start to end of input".
	Ranges [][2]int
}

// Linenupper filters lines by specified line number ranges.
type Linenupper struct {
	ranges [][2]int
}

// New creates a Linenupper from the given Options.
// Returns an error if no ranges are provided or any range is invalid.
func New(opts Options) (*Linenupper, error) {
	if err := validateOptions(opts); err != nil {
		return nil, err
	}
	return &Linenupper{ranges: opts.Ranges}, nil
}

// Apply filters the input lines, returning only those whose 1-indexed
// position falls within at least one of the configured ranges.
func (l *Linenupper) Apply(lines []string) []string {
	var out []string
	total := len(lines)
	for i, line := range lines {
		lineNo := i + 1
		if l.inRange(lineNo, total) {
			out = append(out, line)
		}
	}
	return out
}

func (l *Linenupper) inRange(lineNo, total int) bool {
	for _, r := range l.ranges {
		start := r[0]
		end := r[1]
		if end == 0 {
			end = total
		}
		if lineNo >= start && lineNo <= end {
			return true
		}
	}
	return false
}

// FormatSummary returns a human-readable description of the active ranges.
func FormatSummary(opts Options) string {
	if len(opts.Ranges) == 0 {
		return "line ranges: none"
	}
	s := "line ranges:"
	for _, r := range opts.Ranges {
		if r[1] == 0 {
			s += fmt.Sprintf(" %d-EOF", r[0])
		} else {
			s += fmt.Sprintf(" %d-%d", r[0], r[1])
		}
	}
	return s
}

// validateOptions checks that at least one range is provided and that each
// range has a positive start, a non-negative end, and that end is either 0
// (meaning EOF) or greater than or equal to start.
func validateOptions(opts Options) error {
	if len(opts.Ranges) == 0 {
		return fmt.Errorf("linenum: at least one range must be specified")
	}
	for _, r := range opts.Ranges {
		start, end := r[0], r[1]
		if start < 1 {
			return fmt.Errorf("linenum: range start must be >= 1, got %d", start)
		}
		if end != 0 && end < start {
			return fmt.Errorf("linenum: range end (%d) must be >= start (%d)", end, start)
		}
	}
	return nil
}
