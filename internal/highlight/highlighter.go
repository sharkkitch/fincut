// Package highlight provides regex-based term highlighting for log lines.
package highlight

import (
	"fmt"
	"regexp"
	"strings"
)

const (
	ansiBold   = "\033[1m"
	ansiYellow = "\033[33m"
	ansiRed    = "\033[31m"
	ansiCyan   = "\033[36m"
	ansiReset  = "\033[0m"
)

var colorCycle = []string{ansiYellow, ansiCyan, ansiRed}

// Highlighter applies regex-based highlights to text lines.
type Highlighter struct {
	patterns []*regexp.Regexp
	bold     bool
}

// Options configures a Highlighter.
type Options struct {
	Patterns []string
	Bold     bool
}

// New creates a Highlighter from the given options.
// Returns an error if any pattern fails to compile.
func New(opts Options) (*Highlighter, error) {
	if len(opts.Patterns) == 0 {
		return nil, fmt.Errorf("highlight: at least one pattern required")
	}
	compiled := make([]*regexp.Regexp, 0, len(opts.Patterns))
	for _, p := range opts.Patterns {
		re, err := regexp.Compile(p)
		if err != nil {
			return nil, fmt.Errorf("highlight: invalid pattern %q: %w", p, err)
		}
		compiled = append(compiled, re)
	}
	return &Highlighter{patterns: compiled, bold: opts.Bold}, nil
}

// Apply returns the line with all matching substrings wrapped in ANSI color codes.
func (h *Highlighter) Apply(line string) string {
	result := line
	for i, re := range h.patterns {
		color := colorCycle[i%len(colorCycle)]
		prefix := color
		if h.bold {
			prefix = ansiBold + color
		}
		result = re.ReplaceAllStringFunc(result, func(match string) string {
			return prefix + match + ansiReset
		})
	}
	return result
}

// ApplyAll applies highlighting to every line in the slice.
func (h *Highlighter) ApplyAll(lines []string) []string {
	out := make([]string, len(lines))
	for i, l := range lines {
		out[i] = h.Apply(l)
	}
	return out
}

// StripANSI removes ANSI escape sequences from a string.
func StripANSI(s string) string {
	ansiEscape := regexp.MustCompile(`\033\[[0-9;]*m`)
	return strings.TrimRight(ansiEscape.ReplaceAllString(s, ""), " ")
}
