package clamp

import (
	"fmt"
	"regexp"
)

// Clamper restricts output to lines whose numeric capture group falls within [Min, Max].
type Clamper struct {
	re  *regexp.Regexp
	min float64
	max float64
}

// Options configures the Clamper.
type Options struct {
	Pattern string
	Min     *float64
	Max     *float64
}

// New returns a Clamper or an error if options are invalid.
func New(opts Options) (*Clamper, error) {
	if err := validateOptions(opts); err != nil {
		return nil, err
	}
	re, err := regexp.Compile(opts.Pattern)
	if err != nil {
		return nil, fmt.Errorf("clamp: invalid pattern: %w", err)
	}
	if re.NumSubexp() < 1 {
		return nil, fmt.Errorf("clamp: pattern must contain at least one capture group")
	}
	var mn, mx float64
	if opts.Min != nil {
		mn = *opts.Min
	}
	if opts.Max != nil {
		mx = *opts.Max
	} else {
		mx = 1<<53 - 1
	}
	return &Clamper{re: re, min: mn, max: mx}, nil
}

// Apply filters lines, keeping only those whose captured number is within [min, max].
func (c *Clamper) Apply(lines []string) []string {
	out := make([]string, 0, len(lines))
	for _, l := range lines {
		m := c.re.FindStringSubmatch(l)
		if m == nil {
			continue
		}
		var v float64
		if _, err := fmt.Sscanf(m[1], "%g", &v); err != nil {
			continue
		}
		if v >= c.min && v <= c.max {
			out = append(out, l)
		}
	}
	return out
}

// FormatSummary returns a one-line description of the clamp range.
func (c *Clamper) FormatSummary(matched, total int) string {
	return fmt.Sprintf("clamp: %d/%d lines within [%g, %g]", matched, total, c.min, c.max)
}
