// Package offset provides line-number and byte-offset extraction
// from structured log lines, enabling precise range queries.
package offset

import (
	"errors"
	"fmt"
)

// Entry holds the result of an offset computation for a single line.
type Entry struct {
	Line       int
	ByteStart  int64
	ByteEnd    int64
	Content    string
}

// Options configures the Offsetter.
type Options struct {
	// StartLine is the first line to include (1-based, inclusive).
	StartLine int
	// EndLine is the last line to include (1-based, inclusive). 0 means no limit.
	EndLine int
}

// Offsetter computes byte offsets for lines within a range.
type Offsetter struct {
	opts Options
}

// New creates a new Offsetter with the given options.
func New(opts Options) (*Offsetter, error) {
	if opts.StartLine < 1 {
		return nil, errors.New("offset: StartLine must be >= 1")
	}
	if opts.EndLine != 0 && opts.EndLine < opts.StartLine {
		return nil, fmt.Errorf("offset: EndLine (%d) must be >= StartLine (%d)", opts.EndLine, opts.StartLine)
	}
	return &Offsetter{opts: opts}, nil
}

// Apply processes lines and returns offset entries for those within range.
// Each line is assumed to be terminated by a newline (1 byte) for byte accounting.
func (o *Offsetter) Apply(lines []string) []Entry {
	var (
		results  []Entry
		bytePos  int64
	)
	for i, line := range lines {
		lineNum := i + 1
		lineLen := int64(len(line)) + 1 // +1 for newline
		start := bytePos
		end := bytePos + lineLen - 1
		bytePos += lineLen

		if lineNum < o.opts.StartLine {
			continue
		}
		if o.opts.EndLine != 0 && lineNum > o.opts.EndLine {
			break
		}
		results = append(results, Entry{
			Line:      lineNum,
			ByteStart: start,
			ByteEnd:   end,
			Content:   line,
		})
	}
	return results
}
