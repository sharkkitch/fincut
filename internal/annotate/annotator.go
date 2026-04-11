// Package annotate provides line annotation functionality, prepending
// metadata such as line numbers, timestamps, or custom prefixes to each
// line of structured log output.
package annotate

import (
	"fmt"
	"strings"
	"time"
)

// Annotator prepends configurable metadata to each log line.
type Annotator struct {
	opts Options
}

// New creates a new Annotator with the given options.
// Returns an error if the options are invalid.
func New(opts Options) (*Annotator, error) {
	if err := validateOptions(opts); err != nil {
		return nil, err
	}
	return &Annotator{opts: opts}, nil
}

// Apply annotates each line in the input slice and returns the result.
func (a *Annotator) Apply(lines []string) []string {
	out := make([]string, 0, len(lines))
	for i, line := range lines {
		out = append(out, a.annotate(i+1, line))
	}
	return out
}

func (a *Annotator) annotate(lineNum int, line string) string {
	var parts []string

	if a.opts.LineNumbers {
		parts = append(parts, fmt.Sprintf("%d", lineNum))
	}

	if a.opts.Timestamp {
		ts := time.Now().UTC().Format(a.opts.TimestampFormat)
		parts = append(parts, ts)
	}

	if a.opts.Prefix != "" {
		parts = append(parts, a.opts.Prefix)
	}

	if len(parts) == 0 {
		return line
	}

	return strings.Join(parts, a.opts.Separator) + a.opts.Separator + line
}
