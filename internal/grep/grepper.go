// Package grep provides line-level pattern matching with context window support.
package grep

import (
	"fmt"
	"regexp"
)

// Options configures the Grepper behaviour.
type Options struct {
	Patterns      []string
	ContextBefore int
	ContextAfter  int
	Invert        bool
}

// Match represents a matched line together with its surrounding context.
type Match struct {
	LineNumber int
	Line       string
	Before     []string
	After      []string
}

// Grepper searches lines for one or more patterns.
type Grepper struct {
	regs    []*regexp.Regexp
	opts    Options
}

// New creates a Grepper from the given Options, compiling all patterns.
func New(opts Options) (*Grepper, error) {
	if len(opts.Patterns) == 0 {
		return nil, fmt.Errorf("grep: at least one pattern is required")
	}
	if opts.ContextBefore < 0 || opts.ContextAfter < 0 {
		return nil, fmt.Errorf("grep: context window must be non-negative")
	}
	regs := make([]*regexp.Regexp, 0, len(opts.Patterns))
	for _, p := range opts.Patterns {
		re, err := regexp.Compile(p)
		if err != nil {
			return nil, fmt.Errorf("grep: invalid pattern %q: %w", p, err)
		}
		regs = append(regs, re)
	}
	return &Grepper{regs: regs, opts: opts}, nil
}

// Apply scans lines and returns all matches with context.
func (g *Grepper) Apply(lines []string) []Match {
	var matches []Match
	for i, line := range lines {
		matched := g.matchesAny(line)
		if g.opts.Invert {
			matched = !matched
		}
		if !matched {
			continue
		}
		m := Match{
			LineNumber: i + 1,
			Line:       line,
			Before:     contextSlice(lines, i-g.opts.ContextBefore, i),
			After:      contextSlice(lines, i+1, i+1+g.opts.ContextAfter),
		}
		matches = append(matches, m)
	}
	return matches
}

func (g *Grepper) matchesAny(line string) bool {
	for _, re := range g.regs {
		if re.MatchString(line) {
			return true
		}
	}
	return false
}

func contextSlice(lines []string, from, to int) []string {
	if from < 0 {
		from = 0
	}
	if to > len(lines) {
		to = len(lines)
	}
	if from >= to {
		return nil
	}
	result := make([]string, to-from)
	copy(result, lines[from:to])
	return result
}
