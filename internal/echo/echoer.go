// Package echo provides a pass-through processor that optionally prints
// each line to a secondary writer before forwarding it unchanged. Useful
// for debugging pipelines or teeing output to a log.
package echo

import (
	"fmt"
	"io"
)

// Echoer forwards every line unchanged while optionally writing a copy to
// a secondary writer.
type Echoer struct {
	writer  io.Writer
	prefix  string
	counter int
}

// Options configures the Echoer.
type Options struct {
	// Writer is the secondary destination for echoed lines. If nil, echoing
	// is disabled and Apply acts as a pure pass-through.
	Writer io.Writer

	// Prefix is an optional string prepended to each echoed line.
	Prefix string
}

// New creates an Echoer from the provided Options.
// It returns an error if Writer is nil (no-op echoers should simply not be
// used, but we enforce explicitness).
func New(opts Options) (*Echoer, error) {
	if err := validateOptions(opts); err != nil {
		return nil, err
	}
	return &Echoer{
		writer: opts.Writer,
		prefix: opts.Prefix,
	}, nil
}

// Apply echoes every line to the secondary writer (if set) and returns the
// original slice unmodified.
func (e *Echoer) Apply(lines []string) []string {
	if e.writer == nil {
		return lines
	}
	for _, line := range lines {
		e.counter++
		if e.prefix != "" {
			fmt.Fprintf(e.writer, "%s%s\n", e.prefix, line)
		} else {
			fmt.Fprintln(e.writer, line)
		}
	}
	return lines
}

// FormatSummary returns a human-readable summary of how many lines were echoed.
func (e *Echoer) FormatSummary() string {
	return fmt.Sprintf("echoed %d line(s)", e.counter)
}
