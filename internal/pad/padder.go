// Package pad provides line padding to a fixed width using configurable fill characters.
package pad

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

// Options configures the Padder.
type Options struct {
	// Width is the target rune width for each line. Required, must be > 0.
	Width int
	// Fill is the character used to pad. Defaults to a single space.
	Fill rune
	// Left aligns padding to the left (right-aligns text) instead of the default right padding.
	Left bool
}

// Padder pads each input line to a fixed rune width.
type Padder struct {
	opts Options
}

// New creates a Padder from opts, returning an error if opts are invalid.
func New(opts Options) (*Padder, error) {
	if err := validateOptions(&opts); err != nil {
		return nil, err
	}
	return &Padder{opts: opts}, nil
}

// Apply pads each line in lines to the configured width and returns the result.
func (p *Padder) Apply(lines []string) []string {
	out := make([]string, len(lines))
	for i, line := range lines {
		out[i] = p.pad(line)
	}
	return out
}

func (p *Padder) pad(line string) string {
	n := utf8.RuneCountInString(line)
	if n >= p.opts.Width {
		return line
	}
	fill := strings.Repeat(string(p.opts.Fill), p.opts.Width-n)
	if p.opts.Left {
		return fill + line
	}
	return line + fill
}

// FormatSummary returns a human-readable summary of padding settings.
func FormatSummary(opts Options) string {
	align := "right"
	if opts.Left {
		align = "left"
	}
	return fmt.Sprintf("pad width=%d fill=%q align=%s", opts.Width, opts.Fill, align)
}
