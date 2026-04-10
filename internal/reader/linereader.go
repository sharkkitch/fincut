// Package reader provides utilities for reading and parsing structured log files.
package reader

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

// LineReader reads lines from a source with optional byte-range trimming.
type LineReader struct {
	start int64
	end   int64
}

// NewLineReader creates a LineReader with optional start/end byte offsets.
// Pass -1 for end to read until EOF.
func NewLineReader(start, end int64) (*LineReader, error) {
	if start < 0 {
		return nil, fmt.Errorf("start offset must be non-negative, got %d", start)
	}
	if end != -1 && end < start {
		return nil, fmt.Errorf("end offset (%d) must be greater than start offset (%d)", end, start)
	}
	return &LineReader{start: start, end: end}, nil
}

// ReadFile opens the named file and returns lines within the configured byte range.
func (lr *LineReader) ReadFile(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("opening file %q: %w", path, err)
	}
	defer f.Close()
	return lr.ReadFrom(f)
}

// ReadFrom reads lines from r, respecting the configured byte-range window.
func (lr *LineReader) ReadFrom(r io.Reader) ([]string, error) {
	if lr.start > 0 {
		if seeker, ok := r.(io.Seeker); ok {
			if _, err := seeker.Seek(lr.start, io.SeekStart); err != nil {
				return nil, fmt.Errorf("seeking to offset %d: %w", lr.start, err)
			}
		} else {
			if _, err := io.CopyN(io.Discard, r, lr.start); err != nil {
				return nil, fmt.Errorf("skipping %d bytes: %w", lr.start, err)
			}
		}
	}

	var limited io.Reader = r
	if lr.end != -1 {
		limited = io.LimitReader(r, lr.end-lr.start)
	}

	var lines []string
	scanner := bufio.NewScanner(limited)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scanning lines: %w", err)
	}
	return lines, nil
}
