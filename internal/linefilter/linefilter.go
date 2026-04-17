// Package linefilter provides line-level include/exclude filtering by regex.
package linefilter

import (
	"fmt"
	"regexp"
)

// Options configures the LineFilter.
type Options struct {
	// Include is a list of regex patterns; a line must match at least one.
	Include []string
	// Exclude is a list of regex patterns; a line matching any is dropped.
	Exclude []string
}

// LineFilter applies include and exclude regex rules to lines.
type LineFilter struct {
	include []*regexp.Regexp
	exclude []*regexp.Regexp
}

// New constructs a LineFilter from Options.
func New(opts Options) (*LineFilter, error) {
	if len(opts.Include) == 0 && len(opts.Exclude) == 0 {
		return nil, fmt.Errorf("linefilter: at least one include or exclude pattern required")
	}
	lf := &LineFilter{}
	for _, p := range opts.Include {
		re, err := regexp.Compile(p)
		if err != nil {
			return nil, fmt.Errorf("linefilter: invalid include pattern %q: %w", p, err)
		}
		lf.include = append(lf.include, re)
	}
	for _, p := range opts.Exclude {
		re, err := regexp.Compile(p)
		if err != nil {
			return nil, fmt.Errorf("linefilter: invalid exclude pattern %q: %w", p, err)
		}
		lf.exclude = append(lf.exclude, re)
	}
	return lf, nil
}

// Apply filters lines, returning only those that pass include/exclude rules.
func (lf *LineFilter) Apply(lines []string) []string {
	out := make([]string, 0, len(lines))
	for _, line := range lines {
		if lf.excluded(line) {
			continue
		}
		if len(lf.include) > 0 && !lf.included(line) {
			continue
		}
		out = append(out, line)
	}
	return out
}

func (lf *LineFilter) included(line string) bool {
	for _, re := range lf.include {
		if re.MatchString(line) {
			return true
		}
	}
	return false
}

func (lf *LineFilter) excluded(line string) bool {
	for _, re := range lf.exclude {
		if re.MatchString(line) {
			return true
		}
	}
	return false
}

// FormatSummary returns a human-readable summary of filter results.
func FormatSummary(in, out int) string {
	dropped := in - out
	return fmt.Sprintf("linefilter: %d in, %d out, %d dropped", in, out, dropped)
}
