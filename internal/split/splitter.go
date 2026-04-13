package split

import (
	"fmt"
	"regexp"
)

// Splitter divides a stream of lines into named groups based on a delimiter
// pattern. Each time the pattern matches, a new group begins.
type Splitter struct {
	pattern   *regexp.Regexp
	keepDelim bool
	label     string
}

// Group holds a named collection of lines produced by the Splitter.
type Group struct {
	Label string
	Lines []string
}

// Options configures a Splitter.
type Options struct {
	// Pattern is the regex that triggers a new group boundary.
	Pattern string
	// KeepDelimiter, when true, includes the matching line in the new group.
	KeepDelimiter bool
	// Label is the prefix used when auto-generating group names.
	Label string
}

// New constructs a Splitter from the given options.
func New(opts Options) (*Splitter, error) {
	if err := validateOptions(opts); err != nil {
		return nil, err
	}
	re, err := regexp.Compile(opts.Pattern)
	if err != nil {
		return nil, fmt.Errorf("split: invalid pattern %q: %w", opts.Pattern, err)
	}
	label := opts.Label
	if label == "" {
		label = "group"
	}
	return &Splitter{pattern: re, keepDelim: opts.KeepDelimiter, label: label}, nil
}

// Apply splits lines into groups, returning at least one group even when no
// delimiter is ever matched.
func (s *Splitter) Apply(lines []string) []Group {
	var groups []Group
	current := Group{Label: fmt.Sprintf("%s-1", s.label)}

	for _, line := range lines {
		if s.pattern.MatchString(line) {
			groups = append(groups, current)
			current = Group{
				Label: fmt.Sprintf("%s-%d", s.label, len(groups)+1),
			}
			if s.keepDelim {
				current.Lines = append(current.Lines, line)
			}
			continue
		}
		current.Lines = append(current.Lines, line)
	}
	groups = append(groups, current)
	return groups
}
