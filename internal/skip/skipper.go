// Package skip provides a line skipper that drops every Nth line from input.
package skip

import "fmt"

// Options configures the Skipper.
type Options struct {
	// Every drops every Nth line (1-based). Must be >= 2.
	Every int
	// Offset shifts the skip window by this many lines before starting.
	Offset int
}

// Skipper drops every Nth line from a slice of lines.
type Skipper struct {
	opts Options
}

// New creates a Skipper with the given options.
func New(opts Options) (*Skipper, error) {
	if err := validateOptions(opts); err != nil {
		return nil, err
	}
	return &Skipper{opts: opts}, nil
}

func validateOptions(opts Options) error {
	if opts.Every < 2 {
		return fmt.Errorf("skip: Every must be >= 2, got %d", opts.Every)
	}
	if opts.Offset < 0 {
		return fmt.Errorf("skip: Offset must be >= 0, got %d", opts.Offset)
	}
	return nil
}

// Apply returns a new slice with every Nth line removed.
// Line numbering starts at 1 after the offset is applied.
func (s *Skipper) Apply(lines []string) []string {
	if len(lines) == 0 {
		return lines
	}
	out := make([]string, 0, len(lines))
	for i, line := range lines {
		// Adjust index by offset; position is 1-based within the effective window.
		pos := (i - s.opts.Offset) + 1
		if pos > 0 && pos%s.opts.Every == 0 {
			continue
		}
		out = append(out, line)
	}
	return out
}

// FormatSummary returns a human-readable summary of the skip operation.
func FormatSummary(in, out int) string {
	return fmt.Sprintf("skip: %d in, %d out, %d dropped", in, out, in-out)
}
