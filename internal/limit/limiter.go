// Package limit provides a line-count and byte-count hard limiter that stops
// processing once a configured ceiling is reached.
package limit

import "fmt"

// Options controls the behaviour of the Limiter.
type Options struct {
	// MaxLines stops output after this many lines (0 = unlimited).
	MaxLines int
	// MaxBytes stops output after this many bytes (0 = unlimited).
	MaxBytes int64
}

// Limiter enforces hard ceilings on line count and total byte volume.
type Limiter struct {
	opts Options
}

// New returns a Limiter configured with opts.
// It returns an error if both MaxLines and MaxBytes are zero, or if either
// value is negative.
func New(opts Options) (*Limiter, error) {
	if err := validateOptions(opts); err != nil {
		return nil, err
	}
	return &Limiter{opts: opts}, nil
}

// Apply processes lines and returns those that fall within the configured
// ceiling.  The returned slice is always non-nil.
func (l *Limiter) Apply(lines []string) []string {
	out := make([]string, 0, len(lines))
	var bytesSeen int64
	for i, line := range lines {
		if l.opts.MaxLines > 0 && i >= l.opts.MaxLines {
			break
		}
		lineBytes := int64(len(line))
		if l.opts.MaxBytes > 0 && bytesSeen+lineBytes > l.opts.MaxBytes {
			break
		}
		out = append(out, line)
		bytesSeen += lineBytes
	}
	return out
}

// FormatSummary returns a human-readable description of the active limits.
func (l *Limiter) FormatSummary() string {
	switch {
	case l.opts.MaxLines > 0 && l.opts.MaxBytes > 0:
		return fmt.Sprintf("limit: %d lines / %d bytes", l.opts.MaxLines, l.opts.MaxBytes)
	case l.opts.MaxLines > 0:
		return fmt.Sprintf("limit: %d lines", l.opts.MaxLines)
	default:
		return fmt.Sprintf("limit: %d bytes", l.opts.MaxBytes)
	}
}
