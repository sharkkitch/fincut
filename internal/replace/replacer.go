// Package replace provides line-level string replacement using regex patterns.
package replace

import (
	"fmt"
	"regexp"
)

// Replacement pairs a compiled pattern with its substitution string.
type Replacement struct {
	re          *regexp.Regexp
	substitution string
}

// Replacer applies one or more regex substitutions to each line.
type Replacer struct {
	replacements []Replacement
	literal      bool
}

// Options configures the Replacer.
type Options struct {
	// Patterns is a list of regex=replacement pairs, e.g. ["foo=bar", `\d+=NUM`].
	Patterns []string
	// Literal treats replacement strings as literals (no backreference expansion).
	Literal bool
}

// New creates a Replacer from the given Options.
func New(opts Options) (*Replacer, error) {
	if err := validateOptions(opts); err != nil {
		return nil, err
	}

	replacements := make([]Replacement, 0, len(opts.Patterns))
	for _, p := range opts.Patterns {
		pattern, sub, err := splitPair(p)
		if err != nil {
			return nil, err
		}
		re, err := regexp.Compile(pattern)
		if err != nil {
			return nil, fmt.Errorf("replace: invalid pattern %q: %w", pattern, err)
		}
		replacements = append(replacements, Replacement{re: re, substitution: sub})
	}

	return &Replacer{replacements: replacements, literal: opts.Literal}, nil
}

// Apply runs all replacements against each line and returns the modified lines.
func (r *Replacer) Apply(lines []string) []string {
	out := make([]string, len(lines))
	for i, line := range lines {
		for _, rep := range r.replacements {
			if r.literal {
				line = rep.re.ReplaceAllLiteralString(line, rep.substitution)
			} else {
				line = rep.re.ReplaceAllString(line, rep.substitution)
			}
		}
		out[i] = line
	}
	return out
}

// splitPair splits a "pattern=replacement" string at the first '=' separator.
func splitPair(s string) (pattern, replacement string, err error) {
	for i, c := range s {
		if c == '=' {
			return s[:i], s[i+1:], nil
		}
	}
	return "", "", fmt.Errorf("replace: pair %q missing '=' separator", s)
}
