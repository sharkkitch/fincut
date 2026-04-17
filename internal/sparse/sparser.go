// Package sparse provides line-skipping based on a fixed stride across input.
package sparse

import "fmt"

// Options configures the Sparser.
type Options struct {
	// Every nth line is kept (1 = keep all, 2 = keep every other, etc.).
	Every int
	// Offset is the zero-based index of the first line to keep within each stride.
	Offset int
}

// Sparser filters lines by keeping every Nth line.
type Sparser struct {
	every  int
	offset int
}

// New creates a Sparser from opts.
func New(opts Options) (*Sparser, error) {
	if err := validateOptions(opts); err != nil {
		return nil, err
	}
	return &Sparser{every: opts.Every, offset: opts.Offset}, nil
}

// Apply returns the subset of lines matching the stride pattern.
func (s *Sparser) Apply(lines []string) []string {
	if len(lines) == 0 {
		return []string{}
	}
	out := make([]string, 0, len(lines)/s.every+1)
	for i, line := range lines {
		if i%s.every == s.offset {
			out = append(out, line)
		}
	}
	return out
}

// FormatSummary returns a human-readable summary of sparse parameters.
func FormatSummary(total, kept int, opts Options) string {
	return fmt.Sprintf("sparse: kept %d/%d lines (every=%d offset=%d)", kept, total, opts.Every, opts.Offset)
}
