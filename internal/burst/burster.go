// Package burst detects and groups high-frequency line bursts within a
// sliding time window, useful for identifying log storms or rate spikes.
package burst

import (
	"fmt"
	"regexp"
	"time"
)

// Burst represents a detected burst of lines within a time window.
type Burst struct {
	Start int
	End   int
	Lines []string
	Rate  float64 // lines per second
}

// Burster groups consecutive lines that exceed a threshold rate.
type Burster struct {
	opts Options
	tsRe *regexp.Regexp
}

// New returns a Burster or an error if options are invalid.
func New(opts Options) (*Burster, error) {
	if err := validateOptions(&opts); err != nil {
		return nil, err
	}
	re, err := regexp.Compile(opts.TimestampPattern)
	if err != nil {
		return nil, fmt.Errorf("burst: invalid timestamp pattern: %w", err)
	}
	return &Burster{opts: opts, tsRe: re}, nil
}

// Apply scans lines for bursts and returns detected Burst slices.
func (b *Burster) Apply(lines []string) ([]Burst, error) {
	var bursts []Burst
	var window []string
	var windowStart time.Time
	var windowIdx int

	for i, line := range lines {
		ts, ok := b.parseTimestamp(line)
		if !ok {
			window = append(window, line)
			continue
		}

		if windowStart.IsZero() {
			windowStart = ts
			windowIdx = i
		}

		window = append(window, line)
		elapsed := ts.Sub(windowStart)

		if elapsed >= b.opts.Window {
			rate := float64(len(window)) / elapsed.Seconds()
			if rate >= b.opts.Threshold {
				bursts = append(bursts, Burst{
					Start: windowIdx,
					End:   i,
					Lines: append([]string(nil), window...),
					Rate:  rate,
				})
			}
			window = window[:0]
			windowStart = time.Time{}
		}
	}
	return bursts, nil
}

func (b *Burster) parseTimestamp(line string) (time.Time, bool) {
	m := b.tsRe.FindString(line)
	if m == "" {
		return time.Time{}, false
	}
	t, err := time.Parse(b.opts.TimestampLayout, m)
	if err != nil {
		return time.Time{}, false
	}
	return t, true
}
