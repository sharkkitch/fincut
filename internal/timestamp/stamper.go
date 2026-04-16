package timestamp

import (
	"fmt"
	"time"
)

// Stamper prepends or appends a timestamp to each line.
type Stamper struct {
	opts Options
}

// Options configures the Stamper.
type Options struct {
	Format  string // time format layout, defaults to time.RFC3339
	Prepend bool   // prepend timestamp; mutually exclusive with Append
	Append  bool   // append timestamp
	Sep     string // separator between timestamp and line, defaults to " "
}

// New creates a Stamper from opts.
func New(opts Options) (*Stamper, error) {
	if err := validateOptions(&opts); err != nil {
		return nil, err
	}
	return &Stamper{opts: opts}, nil
}

// Apply stamps each line and returns the result.
func (s *Stamper) Apply(lines []string) []string {
	out := make([]string, len(lines))
	ts := time.Now().Format(s.opts.Format)
	for i, l := range lines {
		if s.opts.Prepend {
			out[i] = ts + s.opts.Sep + l
		} else {
			out[i] = l + s.opts.Sep + ts
		}
	}
	return out
}

// FormatSummary returns a human-readable summary.
func FormatSummary(total int, layout string) string {
	return fmt.Sprintf("stamped %d lines with layout %q", total, layout)
}
