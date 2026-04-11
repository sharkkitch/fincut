package mask

import (
	"fmt"
	"regexp"
	"strings"
)

// Masker replaces sensitive patterns in log lines with a redaction string.
type Masker struct {
	patterns    []*regexp.Regexp
	replacement string
}

// Options configures the Masker.
type Options struct {
	// Patterns is a list of regex strings whose matches will be redacted.
	Patterns []string
	// Replacement is the string used in place of matched text. Defaults to "[REDACTED]".
	Replacement string
}

// New creates a Masker from the given options.
func New(opts Options) (*Masker, error) {
	if len(opts.Patterns) == 0 {
		return nil, fmt.Errorf("mask: at least one pattern is required")
	}
	replacement := opts.Replacement
	if replacement == "" {
		replacement = "[REDACTED]"
	}
	compiled := make([]*regexp.Regexp, 0, len(opts.Patterns))
	for _, p := range opts.Patterns {
		re, err := regexp.Compile(p)
		if err != nil {
			return nil, fmt.Errorf("mask: invalid pattern %q: %w", p, err)
		}
		compiled = append(compiled, re)
	}
	return &Masker{patterns: compiled, replacement: replacement}, nil
}

// Apply redacts all pattern matches in the given lines and returns the result.
func (m *Masker) Apply(lines []string) []string {
	out := make([]string, len(lines))
	for i, line := range lines {
		out[i] = m.redact(line)
	}
	return out
}

// redact applies all patterns to a single line.
func (m *Masker) redact(line string) string {
	for _, re := range m.patterns {
		line = re.ReplaceAllString(line, m.replacement)
	}
	return line
}

// Summary returns a brief description of the masker configuration.
func (m *Masker) Summary() string {
	return fmt.Sprintf("masker: %d pattern(s), replacement=%q",
		len(m.patterns), m.replacement)
}

// CountRedacted returns how many lines contain at least one redaction.
func (m *Masker) CountRedacted(lines []string) int {
	count := 0
	for _, line := range lines {
		redacted := m.redact(line)
		if !strings.EqualFold(redacted, line) {
			count++
		}
	}
	return count
}
