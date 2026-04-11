// Package tail provides functionality for reading the last N lines or bytes
// from a stream, similar to the Unix tail command.
package tail

import (
	"bufio"
	"fmt"
	"io"
)

// Options configures the Tailer behaviour.
type Options struct {
	// MaxLines is the maximum number of trailing lines to return.
	// Zero means no line limit.
	MaxLines int
	// MaxBytes is the maximum number of trailing bytes to return.
	// Zero means no byte limit.
	MaxBytes int64
}

// Tailer reads the trailing portion of a line-oriented stream.
type Tailer struct {
	opts Options
}

// NewTailer creates a Tailer with the given options.
// Returns an error if both MaxLines and MaxBytes are zero or negative.
func NewTailer(opts Options) (*Tailer, error) {
	if opts.MaxLines < 0 {
		return nil, fmt.Errorf("tail: MaxLines must be non-negative, got %d", opts.MaxLines)
	}
	if opts.MaxBytes < 0 {
		return nil, fmt.Errorf("tail: MaxBytes must be non-negative, got %d", opts.MaxBytes)
	}
	if opts.MaxLines == 0 && opts.MaxBytes == 0 {
		return nil, fmt.Errorf("tail: at least one of MaxLines or MaxBytes must be set")
	}
	return &Tailer{opts: opts}, nil
}

// Apply reads all lines from r and returns the trailing subset that satisfies
// the configured constraints. Lines are returned in original order.
func (t *Tailer) Apply(r io.Reader) ([]string, error) {
	scanner := bufio.NewScanner(r)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("tail: reading input: %w", err)
	}

	// Apply MaxLines constraint first.
	if t.opts.MaxLines > 0 && len(lines) > t.opts.MaxLines {
		lines = lines[len(lines)-t.opts.MaxLines:]
	}

	// Apply MaxBytes constraint by walking backwards.
	if t.opts.MaxBytes > 0 {
		var total int64
		start := len(lines)
		for i := len(lines) - 1; i >= 0; i-- {
			// +1 accounts for the newline that was stripped by the scanner.
			total += int64(len(lines[i])) + 1
			if total > t.opts.MaxBytes {
				break
			}
			start = i
		}
		lines = lines[start:]
	}

	return lines, nil
}
