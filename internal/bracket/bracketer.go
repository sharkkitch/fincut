// Package bracket wraps matching lines with a configurable open/close string pair.
package bracket

import (
	"fmt"
	"regexp"
)

// Options configures the Bracketer.
type Options struct {
	// Pattern is the regex that selects lines to wrap.
	Pattern string
	// Open is the string prepended to matching lines.
	Open string
	// Close is the string appended to matching lines.
	Close string
	// WrapAll, when true, wraps every line regardless of Pattern.
	WrapAll bool
}

// Bracketer wraps matching lines with Open/Close strings.
type Bracketer struct {
	opts Options
	re   *regexp.Regexp
}

// New returns a Bracketer or an error if the options are invalid.
func New(opts Options) (*Bracketer, error) {
	if err := validateOptions(opts); err != nil {
		return nil, err
	}
	var re *regexp.Regexp
	if !opts.WrapAll {
		var err error
		re, err = regexp.Compile(opts.Pattern)
		if err != nil {
			return nil, fmt.Errorf("bracket: invalid pattern: %w", err)
		}
	}
	return &Bracketer{opts: opts, re: re}, nil
}

// Apply wraps lines that match the configured pattern (or all lines if WrapAll
// is set) with the Open and Close strings. The original slice is not modified.
func (b *Bracketer) Apply(lines []string) []string {
	out := make([]string, len(lines))
	for i, l := range lines {
		if b.opts.WrapAll || b.re.MatchString(l) {
			out[i] = b.opts.Open + l + b.opts.Close
		} else {
			out[i] = l
		}
	}
	return out
}

// FormatSummary returns a one-line summary of how many lines were wrapped.
func FormatSummary(total, wrapped int) string {
	return fmt.Sprintf("bracket: %d/%d lines wrapped", wrapped, total)
}
