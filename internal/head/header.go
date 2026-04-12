// Package head provides functionality for extracting the first N lines
// or first N bytes from a stream of log lines.
package head

import (
	"fmt"
)

// Options configures the Header behaviour.
type Options struct {
	// MaxLines is the maximum number of lines to retain. Zero means no limit.
	MaxLines int
	// MaxBytes is the maximum number of bytes to retain. Zero means no limit.
	MaxBytes int
}

// Header retains only the leading portion of a line slice.
type Header struct {
	opts Options
}

// New returns a new Header or an error if the options are invalid.
func New(opts Options) (*Header, error) {
	if opts.MaxLines < 0 {
		return nil, fmt.Errorf("head: MaxLines must be non-negative, got %d", opts.MaxLines)
	}
	if opts.MaxBytes < 0 {
		return nil, fmt.Errorf("head: MaxBytes must be non-negative, got %d", opts.MaxBytes)
	}
	if opts.MaxLines == 0 && opts.MaxBytes == 0 {
		return nil, fmt.Errorf("head: at least one of MaxLines or MaxBytes must be set")
	}
	return &Header{opts: opts}, nil
}

// Apply returns the leading lines of input according to the configured limits.
// When both MaxLines and MaxBytes are set the stricter limit wins.
func (h *Header) Apply(lines []string) []string {
	var (
		bytesAccum int
		result     []string
	)
	for i, line := range lines {
		if h.opts.MaxLines > 0 && i >= h.opts.MaxLines {
			break
		}
		lineBytes := len(line) + 1 // +1 for newline
		if h.opts.MaxBytes > 0 && bytesAccum+lineBytes > h.opts.MaxBytes {
			break
		}
		bytesAccum += lineBytes
		result = append(result, line)
	}
	return result
}
