package stripe

import "fmt"

// Striper selects every Nth line from input, with an optional offset.
type Striper struct {
	every  int
	offset int
}

// Options configures the Striper.
type Options struct {
	// Every selects one line out of every N lines.
	Every int
	// Offset is the zero-based index within each group to select.
	Offset int
}

// New creates a Striper from opts.
func New(opts Options) (*Striper, error) {
	if err := validateOptions(opts); err != nil {
		return nil, err
	}
	return &Striper{every: opts.Every, offset: opts.Offset}, nil
}

// Apply returns lines at positions (offset, offset+every, offset+2*every, ...).
func (s *Striper) Apply(lines []string) []string {
	out := make([]string, 0, len(lines)/s.every+1)
	for i, line := range lines {
		if i%s.every == s.offset {
			out = append(out, line)
		}
	}
	return out
}

// FormatSummary returns a human-readable summary of the stripe operation.
func FormatSummary(in, out int, every, offset int) string {
	return fmt.Sprintf("stripe: selected %d/%d lines (every=%d offset=%d)", out, in, every, offset)
}
