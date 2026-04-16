package numrange

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// Options configures the numeric range filter.
type Options struct {
	// Pattern is a regex with one capture group matching a number.
	Pattern string
	// Min is the inclusive lower bound (nil means no lower bound).
	Min *float64
	// Max is the inclusive upper bound (nil means no upper bound).
	Max *float64
	// Invert returns lines that do NOT fall within the range.
	Invert bool
}

// Ranger filters lines by extracting a numeric value and testing it against a range.
type Ranger struct {
	opts Options
	re   *regexp.Regexp
}

// New creates a Ranger from opts.
func New(opts Options) (*Ranger, error) {
	if err := validateOptions(opts); err != nil {
		return nil, err
	}
	re, err := regexp.Compile(opts.Pattern)
	if err != nil {
		return nil, fmt.Errorf("numrange: invalid pattern: %w", err)
	}
	if re.NumSubexp() < 1 {
		return nil, fmt.Errorf("numrange: pattern must contain at least one capture group")
	}
	return &Ranger{opts: opts, re: re}, nil
}

// Apply filters lines, returning only those whose extracted number is within [Min, Max].
func (r *Ranger) Apply(lines []string) []string {
	out := make([]string, 0, len(lines))
	for _, line := range lines {
		m := r.re.FindStringSubmatch(line)
		if m == nil {
			continue
		}
		v, err := strconv.ParseFloat(strings.TrimSpace(m[1]), 64)
		if err != nil {
			continue
		}
		inRange := true
		if r.opts.Min != nil && v < *r.opts.Min {
			inRange = false
		}
		if r.opts.Max != nil && v > *r.opts.Max {
			inRange = false
		}
		if inRange != r.opts.Invert {
			out = append(out, line)
		}
	}
	return out
}

// FormatSummary returns a human-readable summary of filter results.
func FormatSummary(total, kept int) string {
	return fmt.Sprintf("numrange: %d/%d lines kept", kept, total)
}
