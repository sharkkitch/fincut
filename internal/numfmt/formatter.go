// Package numfmt formats numeric fields within log lines.
package numfmt

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// Options configures the numeric formatter.
type Options struct {
	// Pattern is a regex with one capture group matching a numeric value.
	Pattern string
	// Precision is the number of decimal places for float formatting.
	Precision int
	// Thousands inserts comma separators for integers when true.
	Thousands bool
}

// Formatter applies numeric formatting to matched fields in lines.
type Formatter struct {
	re        *regexp.Regexp
	precision int
	thousands bool
}

// New creates a Formatter from the given Options.
func New(opts Options) (*Formatter, error) {
	if err := validateOptions(opts); err != nil {
		return nil, err
	}
	re, err := regexp.Compile(opts.Pattern)
	if err != nil {
		return nil, fmt.Errorf("numfmt: invalid pattern: %w", err)
	}
	if re.NumSubexp() < 1 {
		return nil, fmt.Errorf("numfmt: pattern must contain at least one capture group")
	}
	return &Formatter{re: re, precision: opts.Precision, thousands: opts.Thousands}, nil
}

// Apply formats numeric capture groups in each line and returns the results.
func (f *Formatter) Apply(lines []string) []string {
	out := make([]string, len(lines))
	for i, line := range lines {
		out[i] = f.re.ReplaceAllStringFunc(line, func(match string) string {
			sub := f.re.FindStringSubmatch(match)
			if len(sub) < 2 {
				return match
			}
			raw := sub[1]
			if strings.Contains(raw, ".") {
				v, err := strconv.ParseFloat(raw, 64)
				if err != nil {
					return match
				}
				formatted := strconv.FormatFloat(v, 'f', f.precision, 64)
				return strings.Replace(match, raw, formatted, 1)
			}
			v, err := strconv.ParseInt(raw, 10, 64)
			if err != nil {
				return match
			}
			formatted := formatInt(v, f.thousands)
			return strings.Replace(match, raw, formatted, 1)
		})
	}
	return out
}

func formatInt(n int64, thousands bool) string {
	s := strconv.FormatInt(n, 10)
	if !thousands {
		return s
	}
	neg := ""
	if n < 0 {
		neg = "-"
		s = s[1:]
	}
	var b strings.Builder
	start := len(s) % 3
	if start > 0 {
		b.WriteString(s[:start])
	}
	for i := start; i < len(s); i += 3 {
		if i > 0 || start > 0 {
			b.WriteByte(',')
		}
		b.WriteString(s[i : i+3])
	}
	return neg + b.String()
}
