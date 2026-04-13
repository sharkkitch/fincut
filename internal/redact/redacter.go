// Package redact provides line-level redaction by replacing entire lines
// that match one or more patterns with a configurable replacement string.
package redact

import (
	"fmt"
	"regexp"
)

// Options configures the Redacter.
type Options struct {
	// Patterns is a list of regular expressions; any line fully matched by one
	// of these patterns will be replaced.
	Patterns []string

	// Replacement is the string written in place of a redacted line.
	// Defaults to "[REDACTED]" when empty.
	Replacement string

	// MatchPartial, when true, redacts lines that contain a match anywhere
	// (not just lines where the whole line matches).
	MatchPartial bool
}

// Redacter replaces matching lines with a fixed replacement string.
type Redacter struct {
	opts    Options
	regs    []*regexp.Regexp
	total   int
	redacted int
}

// New constructs a Redacter from opts, compiling all patterns.
func New(opts Options) (*Redacter, error) {
	if err := validateOptions(&opts); err != nil {
		return nil, err
	}
	regs := make([]*regexp.Regexp, 0, len(opts.Patterns))
	for _, p := range opts.Patterns {
		re, err := regexp.Compile(p)
		if err != nil {
			return nil, fmt.Errorf("redact: invalid pattern %q: %w", p, err)
		}
		regs = append(regs, re)
	}
	return &Redacter{opts: opts, regs: regs}, nil
}

// Apply processes lines, replacing any that match a pattern.
func (r *Redacter) Apply(lines []string) []string {
	out := make([]string, len(lines))
	for i, line := range lines {
		r.total++
		if r.matches(line) {
			r.redacted++
			out[i] = r.opts.Replacement
		} else {
			out[i] = line
		}
	}
	return out
}

// Stats returns total and redacted line counts.
func (r *Redacter) Stats() (total, redacted int) {
	return r.total, r.redacted
}

func (r *Redacter) matches(line string) bool {
	for _, re := range r.regs {
		if r.opts.MatchPartial {
			if re.MatchString(line) {
				return true
			}
		} else {
			if re.MatchString(line) && re.FindString(line) == line {
				return true
			}
		}
	}
	return false
}
