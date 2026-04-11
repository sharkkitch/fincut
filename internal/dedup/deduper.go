package dedup

import (
	"crypto/sha256"
	"fmt"
)

// Options configures the Deduper.
type Options struct {
	// WindowSize limits how many recent hashes are tracked (0 = unlimited).
	WindowSize int
	// CaseSensitive controls whether comparison is case-sensitive.
	CaseSensitive bool
}

// Deduper removes duplicate lines from a stream.
type Deduper struct {
	opts   Options
	seen   map[string]struct{}
	window []string
}

// New creates a Deduper with the given options.
func New(opts Options) (*Deduper, error) {
	if opts.WindowSize < 0 {
		return nil, fmt.Errorf("dedup: window size must be non-negative, got %d", opts.WindowSize)
	}
	return &Deduper{
		opts: opts,
		seen: make(map[string]struct{}),
	}, nil
}

// Apply filters duplicate lines and returns only unique ones.
func (d *Deduper) Apply(lines []string) []string {
	out := make([]string, 0, len(lines))
	for _, line := range lines {
		key := d.hash(line)
		if _, exists := d.seen[key]; exists {
			continue
		}
		d.seen[key] = struct{}{}
		if d.opts.WindowSize > 0 {
			d.window = append(d.window, key)
			if len(d.window) > d.opts.WindowSize {
				evict := d.window[0]
				d.window = d.window[1:]
				delete(d.seen, evict)
			}
		}
		out = append(out, line)
	}
	return out
}

// Reset clears the internal deduplication state.
func (d *Deduper) Reset() {
	d.seen = make(map[string]struct{})
	d.window = d.window[:0]
}

func (d *Deduper) hash(line string) string {
	if !d.opts.CaseSensitive {
		line = toLower(line)
	}
	sum := sha256.Sum256([]byte(line))
	return fmt.Sprintf("%x", sum)
}

func toLower(s string) string {
	b := make([]byte, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= 'A' && c <= 'Z' {
			c += 32
		}
		b[i] = c
	}
	return string(b)
}
