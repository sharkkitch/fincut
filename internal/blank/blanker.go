// Package blank provides a processor that inserts blank lines between
// output lines at a configurable interval, useful for improving readability
// of dense log output.
package blank

import "fmt"

// Options configures the Blanker processor.
type Options struct {
	// Every inserts a blank line after every N lines. Must be >= 2.
	Every int
	// Offset is the number of initial lines to skip before inserting blanks.
	Offset int
}

// Blanker inserts blank lines into a slice of strings at regular intervals.
type Blanker struct {
	every  int
	offset int
}

// New creates a new Blanker with the given options.
func New(opts Options) (*Blanker, error) {
	if err := validateOptions(opts); err != nil {
		return nil, err
	}
	return &Blanker{
		every:  opts.Every,
		offset: opts.Offset,
	}, nil
}

// Apply inserts blank lines into lines at the configured interval.
// A blank line is inserted after every N-th line (starting after Offset lines).
func (b *Blanker) Apply(lines []string) []string {
	if len(lines) == 0 {
		return lines
	}
	out := make([]string, 0, len(lines)+len(lines)/b.every+1)
	count := 0
	for i, line := range lines {
		out = append(out, line)
		if i < b.offset {
			continue
		}
		count++
		if count%b.every == 0 {
			out = append(out, "")
		}
	}
	return out
}

// FormatSummary returns a human-readable summary of the blanker configuration.
func FormatSummary(b *Blanker, total, inserted int) string {
	return fmt.Sprintf("blank: inserted %d blank lines into %d input lines (every=%d, offset=%d)",
		inserted, total, b.every, b.offset)
}
