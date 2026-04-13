package prefix

import "fmt"

// Prefixer prepends a fixed string or line-number-based label to each line.
type Prefixer struct {
	opts Options
}

// Options configures the Prefixer.
type Options struct {
	// Text is a literal string prepended to every line.
	Text string
	// LineNumbers, when true, prepends a formatted line number instead of Text.
	LineNumbers bool
	// Separator is placed between the prefix and the original line.
	// Defaults to ": ".
	Separator string
	// Width zero-pads line numbers to this width (only used when LineNumbers is true).
	Width int
}

// New returns a Prefixer or an error if the options are invalid.
func New(opts Options) (*Prefixer, error) {
	if err := validateOptions(&opts); err != nil {
		return nil, err
	}
	return &Prefixer{opts: opts}, nil
}

// Apply prepends the configured prefix to every line and returns the result.
func (p *Prefixer) Apply(lines []string) []string {
	out := make([]string, len(lines))
	for i, line := range lines {
		out[i] = p.prefix(i+1) + p.opts.Separator + line
	}
	return out
}

func (p *Prefixer) prefix(n int) string {
	if p.opts.LineNumbers {
		if p.opts.Width > 0 {
			return fmt.Sprintf("%0*d", p.opts.Width, n)
		}
		return fmt.Sprintf("%d", n)
	}
	return p.opts.Text
}

// FormatSummary returns a one-line description of what was applied.
func FormatSummary(total int, opts Options) string {
	if opts.LineNumbers {
		return fmt.Sprintf("prefixed %d lines with line numbers (sep=%q)", total, opts.Separator)
	}
	return fmt.Sprintf("prefixed %d lines with %q (sep=%q)", total, opts.Text, opts.Separator)
}
