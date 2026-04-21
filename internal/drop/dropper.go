package drop

import (
	"fmt"
	"regexp"
)

// Dropper removes lines matching any of the given patterns.
type Dropper struct {
	patterns []*regexp.Regexp
	invert   bool
}

// Options configures the Dropper.
type Options struct {
	// Patterns is a list of regular expressions; matching lines are dropped.
	Patterns []string
	// Invert reverses the behaviour: only lines that match are kept.
	Invert bool
}

// New creates a Dropper from the provided options.
func New(opts Options) (*Dropper, error) {
	if err := validateOptions(opts); err != nil {
		return nil, err
	}

	compiled := make([]*regexp.Regexp, 0, len(opts.Patterns))
	for _, p := range opts.Patterns {
		re, err := regexp.Compile(p)
		if err != nil {
			return nil, fmt.Errorf("drop: invalid pattern %q: %w", p, err)
		}
		compiled = append(compiled, re)
	}

	return &Dropper{
		patterns: compiled,
		invert:   opts.Invert,
	}, nil
}

// Apply filters lines, returning only those that should be kept.
func (d *Dropper) Apply(lines []string) []string {
	out := make([]string, 0, len(lines))
	for _, line := range lines {
		matched := d.matchesAny(line)
		keep := matched == d.invert // XOR-style: drop matched unless inverted
		if keep {
			out = append(out, line)
		}
	}
	return out
}

// DroppedCount returns how many lines would be dropped from the given slice.
func (d *Dropper) DroppedCount(lines []string) int {
	return len(lines) - len(d.Apply(lines))
}

func (d *Dropper) matchesAny(line string) bool {
	for _, re := range d.patterns {
		if re.MatchString(line) {
			return true
		}
	}
	return false
}

// FormatSummary returns a human-readable summary of the drop operation.
func FormatSummary(total, dropped int) string {
	kept := total - dropped
	return fmt.Sprintf("drop: %d/%d lines kept (%d dropped)", kept, total, dropped)
}
