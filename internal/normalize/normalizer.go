package normalize

import (
	"fmt"
	"strings"
	"unicode"
)

// Normalizer applies whitespace and case normalization to log lines.
type Normalizer struct {
	opts Options
}

// Options configures the Normalizer.
type Options struct {
	// TrimSpace removes leading and trailing whitespace from each line.
	TrimSpace bool
	// CollapseSpaces replaces runs of whitespace with a single space.
	CollapseSpaces bool
	// Lowercase converts each line to lowercase.
	Lowercase bool
	// StripControl removes non-printable control characters.
	StripControl bool
}

// New constructs a Normalizer with the given options.
// Returns an error if no normalization option is enabled.
func New(opts Options) (*Normalizer, error) {
	if err := validateOptions(opts); err != nil {
		return nil, err
	}
	return &Normalizer{opts: opts}, nil
}

// Apply normalizes the given lines according to the configured options.
// Empty lines are preserved unless they become empty after normalization.
func (n *Normalizer) Apply(lines []string) []string {
	out := make([]string, 0, len(lines))
	for _, line := range lines {
		line = n.process(line)
		out = append(out, line)
	}
	return out
}

func (n *Normalizer) process(line string) string {
	if n.opts.StripControl {
		line = stripControl(line)
	}
	if n.opts.CollapseSpaces {
		line = collapseSpaces(line)
	}
	if n.opts.TrimSpace {
		line = strings.TrimSpace(line)
	}
	if n.opts.Lowercase {
		line = strings.ToLower(line)
	}
	return line
}

func collapseSpaces(s string) string {
	var b strings.Builder
	b.Grow(len(s))
	inSpace := false
	for _, r := range s {
		if unicode.IsSpace(r) {
			if !inSpace {
				b.WriteRune(' ')
				inSpace = true
			}
		} else {
			b.WriteRune(r)
			inSpace = false
		}
	}
	return b.String()
}

func stripControl(s string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsControl(r) && r != '\t' && r != '\n' {
			return -1
		}
		return r
	}, s)
}

// Summary returns a human-readable description of active normalizations.
func (n *Normalizer) Summary() string {
	var parts []string
	if n.opts.TrimSpace {
		parts = append(parts, "trim")
	}
	if n.opts.CollapseSpaces {
		parts = append(parts, "collapse-spaces")
	}
	if n.opts.Lowercase {
		parts = append(parts, "lowercase")
	}
	if n.opts.StripControl {
		parts = append(parts, "strip-control")
	}
	return fmt.Sprintf("normalize[%s]", strings.Join(parts, ","))
}
