// Package context provides line-context extraction around matched lines,
// returning a window of surrounding lines for each match.
package context

import (
	"fmt"
	"regexp"
)

// Options configures the Contexter.
type Options struct {
	// Pattern is the regex to match against each line.
	Pattern string
	// Before is the number of lines to include before each match.
	Before int
	// After is the number of lines to include after each match.
	After int
	// Label optionally prefixes each output line with its source line number.
	Label bool
}

// Match holds a matched line and its surrounding context.
type Match struct {
	LineNo int
	Lines  []string
}

// Contexter extracts context windows around matched lines.
type Contexter struct {
	re   *regexp.Regexp
	opts Options
}

// New creates a Contexter from opts, returning an error if opts are invalid.
func New(opts Options) (*Contexter, error) {
	if err := validateOptions(opts); err != nil {
		return nil, err
	}
	re, err := regexp.Compile(opts.Pattern)
	if err != nil {
		return nil, fmt.Errorf("context: invalid pattern: %w", err)
	}
	return &Contexter{re: re, opts: opts}, nil
}

// Apply scans lines and returns a Match for every line matching the pattern.
// Each Match contains up to Before lines preceding and After lines following
// the matched line. Overlapping windows are not merged.
func (c *Contexter) Apply(lines []string) []Match {
	var matches []Match
	for i, line := range lines {
		if !c.re.MatchString(line) {
			continue
		}
		start := i - c.opts.Before
		if start < 0 {
			start = 0
		}
		end := i + c.opts.After + 1
		if end > len(lines) {
			end = len(lines)
		}
		window := make([]string, 0, end-start)
		for j := start; j < end; j++ {
			l := lines[j]
			if c.opts.Label {
				l = fmt.Sprintf("%d: %s", j+1, l)
			}
			window = append(window, l)
		}
		matches = append(matches, Match{LineNo: i + 1, Lines: window})
	}
	return matches
}

func validateOptions(opts Options) error {
	if opts.Pattern == "" {
		return fmt.Errorf("context: pattern must not be empty")
	}
	if opts.Before < 0 {
		return fmt.Errorf("context: before must be >= 0")
	}
	if opts.After < 0 {
		return fmt.Errorf("context: after must be >= 0")
	}
	return nil
}
