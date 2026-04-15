// Package extract provides line extraction by regex capture groups,
// returning only the matched portions of each input line.
package extract

import (
	"fmt"
	"regexp"
)

// Options configures the Extractor.
type Options struct {
	// Pattern is a regular expression with at least one capture group.
	Pattern string
	// Group selects which capture group to extract (1-based). Defaults to 1.
	Group int
	// SkipUnmatched drops lines that do not match instead of passing them through.
	SkipUnmatched bool
}

// Extractor extracts capture-group text from each line.
type Extractor struct {
	re            *regexp.Regexp
	group         int
	skipUnmatched bool
}

// New creates an Extractor from opts.
func New(opts Options) (*Extractor, error) {
	if err := validateOptions(&opts); err != nil {
		return nil, err
	}
	re, err := regexp.Compile(opts.Pattern)
	if err != nil {
		return nil, fmt.Errorf("extract: invalid pattern: %w", err)
	}
	if re.NumSubexp() == 0 {
		return nil, fmt.Errorf("extract: pattern must contain at least one capture group")
	}
	if opts.Group > re.NumSubexp() {
		return nil, fmt.Errorf("extract: group %d exceeds number of capture groups (%d)", opts.Group, re.NumSubexp())
	}
	return &Extractor{
		re:            re,
		group:         opts.Group,
		skipUnmatched: opts.SkipUnmatched,
	}, nil
}

// Apply processes lines and returns extracted capture-group text.
func (e *Extractor) Apply(lines []string) []string {
	out := make([]string, 0, len(lines))
	for _, line := range lines {
		m := e.re.FindStringSubmatch(line)
		if m == nil {
			if !e.skipUnmatched {
				out = append(out, line)
			}
			continue
		}
		out = append(out, m[e.group])
	}
	return out
}

// FormatSummary returns a one-line summary of extraction results.
func FormatSummary(total, extracted, skipped int) string {
	return fmt.Sprintf("extract: %d/%d lines matched, %d skipped", extracted, total, skipped)
}
