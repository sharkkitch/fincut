// Package indent provides line indentation and dedentation for structured log output.
package indent

import (
	"fmt"
	"strings"
)

// Options configures the Indenter.
type Options struct {
	// Depth is the number of indent units to apply. Negative values dedent.
	Depth int
	// Unit is the string used for one level of indentation. Defaults to "  " (two spaces).
	Unit string
	// StripExisting removes leading whitespace before applying indentation.
	StripExisting bool
}

// Indenter applies indentation transformations to lines of text.
type Indenter struct {
	opts Options
	prefix string
}

// New creates a new Indenter with the given options.
// Returns an error if Depth is zero and StripExisting is false (no-op config).
func New(opts Options) (*Indenter, error) {
	if err := validateOptions(&opts); err != nil {
		return nil, err
	}
	prefix := ""
	if opts.Depth > 0 {
		prefix = strings.Repeat(opts.Unit, opts.Depth)
	}
	return &Indenter{opts: opts, prefix: prefix}, nil
}

// Apply transforms the input lines according to the configured indentation.
func (in *Indenter) Apply(lines []string) []string {
	out := make([]string, len(lines))
	for i, line := range lines {
		out[i] = in.transform(line)
	}
	return out
}

func (in *Indenter) transform(line string) string {
	if in.opts.StripExisting {
		line = strings.TrimLeft(line, " \t")
	}
	if in.opts.Depth > 0 {
		return in.prefix + line
	}
	if in.opts.Depth < 0 {
		return dedent(line, in.opts.Unit, -in.opts.Depth)
	}
	return line
}

func dedent(line, unit string, levels int) string {
	for i := 0; i < levels; i++ {
		if strings.HasPrefix(line, unit) {
			line = line[len(unit):]
		} else {
			break
		}
	}
	return line
}

// FormatSummary returns a human-readable summary of the indentation applied.
func FormatSummary(total int, opts Options) string {
	dir := "indented"
	depth := opts.Depth
	if depth < 0 {
		dir = "dedented"
		depth = -depth
	}
	return fmt.Sprintf("%d lines %s by %d level(s) (unit=%q)", total, dir, depth, opts.Unit)
}
