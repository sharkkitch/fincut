// Package center provides a line centerer that pads each line to a fixed width
// by distributing leading and trailing spaces evenly.
package center

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

// Options configures the Centerer.
type Options struct {
	// Width is the total column width to center within. Required; must be > 0.
	Width int
	// Fill is the padding character. Defaults to a single space.
	Fill string
}

// Centerer centers each input line within a fixed column width.
type Centerer struct {
	width int
	fill  string
}

// New creates a Centerer from opts.
func New(opts Options) (*Centerer, error) {
	if err := validateOptions(opts); err != nil {
		return nil, err
	}
	fill := opts.Fill
	if fill == "" {
		fill = " "
	}
	return &Centerer{width: opts.Width, fill: fill}, nil
}

// Apply centers each line in lines and returns the result.
// Lines longer than Width are returned unchanged.
func (c *Centerer) Apply(lines []string) []string {
	out := make([]string, len(lines))
	for i, line := range lines {
		out[i] = c.center(line)
	}
	return out
}

func (c *Centerer) center(line string) string {
	lineWidth := utf8.RuneCountInString(line)
	if lineWidth >= c.width {
		return line
	}
	total := c.width - lineWidth
	left := total / 2
	right := total - left
	return strings.Repeat(c.fill, left) + line + strings.Repeat(c.fill, right)
}

// FormatSummary returns a human-readable summary of the centerer configuration.
func FormatSummary(c *Centerer) string {
	return fmt.Sprintf("center: width=%d fill=%q", c.width, c.fill)
}
