// Package ratelimit provides a rate-limiting filter that drops lines
// exceeding a maximum count within a sliding time window.
package ratelimit

import (
	"fmt"
	"time"
)

// Limiter drops lines that exceed a maximum count per window duration.
type Limiter struct {
	opts    Options
	buckets []bucket
}

type bucket struct {
	start time.Time
	count int
}

// Options configures the rate limiter.
type Options struct {
	// MaxLines is the maximum number of lines allowed per Window.
	MaxLines int
	// Window is the duration of each rate-limit bucket.
	Window time.Duration
	// Now is an optional clock override for testing.
	Now func() time.Time
}

// New creates a Limiter with the given options.
func New(opts Options) (*Limiter, error) {
	if err := validateOptions(opts); err != nil {
		return nil, err
	}
	if opts.Now == nil {
		opts.Now = time.Now
	}
	return &Limiter{opts: opts}, nil
}

// Allow returns true if the line should be passed through, false if it
// should be dropped because the rate limit has been exceeded.
func (l *Limiter) Allow() bool {
	now := l.opts.Now()
	l.evict(now)

	total := 0
	for _, b := range l.buckets {
		total += b.count
	}
	if total >= l.opts.MaxLines {
		return false
	}

	if len(l.buckets) == 0 || now.Sub(l.buckets[len(l.buckets)-1].start) >= l.opts.Window {
		l.buckets = append(l.buckets, bucket{start: now, count: 1})
	} else {
		l.buckets[len(l.buckets)-1].count++
	}
	return true
}

// Apply filters lines, returning only those within the rate limit.
func (l *Limiter) Apply(lines []string) []string {
	out := make([]string, 0, len(lines))
	for _, line := range lines {
		if l.Allow() {
			out = append(out, line)
		}
	}
	return out
}

// FormatSummary returns a human-readable summary of the limiter configuration.
func FormatSummary(maxLines int, window time.Duration) string {
	return fmt.Sprintf("rate-limit: %d lines per %s", maxLines, window)
}

func (l *Limiter) evict(now time.Time) {
	cutoff := now.Add(-l.opts.Window)
	i := 0
	for i < len(l.buckets) && l.buckets[i].start.Before(cutoff) {
		i++
	}
	l.buckets = l.buckets[i:]
}
