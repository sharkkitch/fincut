package between

import (
	"fmt"
	"regexp"
)

// Betweener extracts lines between two regex boundary patterns.
type Betweener struct {
	start   *regexp.Regexp
	end     *regexp.Regexp
	opts    Options
}

// Options configures the Betweener.
type Options struct {
	StartPattern string
	EndPattern   string
	Inclusive    bool // include boundary lines in output
}

// New creates a Betweener from opts.
func New(opts Options) (*Betweener, error) {
	if err := validateOptions(opts); err != nil {
		return nil, err
	}
	start, err := regexp.Compile(opts.StartPattern)
	if err != nil {
		return nil, fmt.Errorf("between: invalid start pattern: %w", err)
	}
	end, err := regexp.Compile(opts.EndPattern)
	if err != nil {
		return nil, fmt.Errorf("between: invalid end pattern: %w", err)
	}
	return &Betweener{start: start, end: end, opts: opts}, nil
}

// Apply returns lines between the first start match and the next end match.
// Multiple regions are supported.
func (b *Betweener) Apply(lines []string) []string {
	var out []string
	inside := false
	for _, line := range lines {
		if !inside && b.start.MatchString(line) {
			inside = true
			if b.opts.Inclusive {
				out = append(out, line)
			}
			continue
		}
		if inside && b.end.MatchString(line) {
			if b.opts.Inclusive {
				out = append(out, line)
			}
			inside = false
			continue
		}
		if inside {
			out = append(out, line)
		}
	}
	return out
}

// FormatSummary returns a human-readable summary string.
func FormatSummary(in, out int) string {
	return fmt.Sprintf("between: %d input lines, %d extracted", in, out)
}
