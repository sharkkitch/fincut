// Package slice provides line-range slicing for log line collections.
// It supports inclusive start/end line numbers (1-based) and optional
// step intervals for sampling every Nth line within the range.
package slice

import (
	"errors"
	"fmt"
)

// Options configures the Slicer behaviour.
type Options struct {
	// Start is the first line to include (1-based, inclusive). Defaults to 1.
	Start int
	// End is the last line to include (1-based, inclusive). 0 means no limit.
	End int
	// Step keeps every Nth line within the range. 0 or 1 means every line.
	Step int
}

// Slicer extracts a sub-range of lines from a slice of strings.
type Slicer struct {
	opts Options
}

// New returns a new Slicer or an error if the options are invalid.
func New(opts Options) (*Slicer, error) {
	if err := validateOptions(&opts); err != nil {
		return nil, err
	}
	return &Slicer{opts: opts}, nil
}

func validateOptions(o *Options) error {
	if o.Start < 1 {
		o.Start = 1
	}
	if o.End != 0 && o.End < o.Start {
		return fmt.Errorf("slice: End (%d) must be >= Start (%d)", o.End, o.Start)
	}
	if o.Step < 0 {
		return errors.New("slice: Step must be >= 0")
	}
	if o.Step == 0 {
		o.Step = 1
	}
	return nil
}

// Apply returns the lines selected by the configured range and step.
// Input lines are treated as 1-based. Lines outside [Start, End] are
// discarded; within the range only lines at positions where
// (pos-Start) % Step == 0 are kept.
func (s *Slicer) Apply(lines []string) []string {
	result := make([]string, 0, len(lines))
	for i, line := range lines {
		pos := i + 1 // convert to 1-based
		if pos < s.opts.Start {
			continue
		}
		if s.opts.End != 0 && pos > s.opts.End {
			break
		}
		if (pos-s.opts.Start)%s.opts.Step == 0 {
			result = append(result, line)
		}
	}
	return result
}
