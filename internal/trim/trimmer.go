// Package trim provides functionality for trimming structured log content
// based on byte ranges and line boundaries.
package trim

import (
	"fmt"
	"io"
	"strings"
)

// Options configures the trimming behavior.
type Options struct {
	MaxLines   int
	MaxBytes   int64
	StripEmpty bool
}

// Trimmer applies trim options to a slice of lines.
type Trimmer struct {
	opts Options
}

// NewTrimmer creates a Trimmer with the given options.
// Returns an error if options are invalid.
func NewTrimmer(opts Options) (*Trimmer, error) {
	if opts.MaxLines < 0 {
		return nil, fmt.Errorf("trim: MaxLines must be non-negative, got %d", opts.MaxLines)
	}
	if opts.MaxBytes < 0 {
		return nil, fmt.Errorf("trim: MaxBytes must be non-negative, got %d", opts.MaxBytes)
	}
	return &Trimmer{opts: opts}, nil
}

// Apply trims the provided lines according to the configured options
// and writes the result to w.
func (t *Trimmer) Apply(lines []string, w io.Writer) (int, error) {
	result := make([]string, 0, len(lines))
	var byteCount int64

	for _, line := range lines {
		if t.opts.StripEmpty && strings.TrimSpace(line) == "" {
			continue
		}
		lineBytes := int64(len(line) + 1) // +1 for newline
		if t.opts.MaxBytes > 0 && byteCount+lineBytes > t.opts.MaxBytes {
			break
		}
		result = append(result, line)
		byteCount += lineBytes
		if t.opts.MaxLines > 0 && len(result) >= t.opts.MaxLines {
			break
		}
	}

	written := 0
	for _, line := range result {
		n, err := fmt.Fprintln(w, line)
		written += n
		if err != nil {
			return written, fmt.Errorf("trim: write error: %w", err)
		}
	}
	return written, nil
}

// Count returns how many lines would survive trimming without writing output.
func (t *Trimmer) Count(lines []string) int {
	var buf strings.Builder
	n, _ := t.Apply(lines, &buf)
	_ = n
	return strings.Count(buf.String(), "\n")
}
