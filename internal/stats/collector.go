package stats

import (
	"fmt"
	"strings"
	"time"
)

// Collector tracks processing statistics for a fincut run.
type Collector struct {
	StartTime     time.Time
	LinesRead     int
	LinesMatched  int
	LinesDropped  int
	BytesRead     int64
	FilterStages  int
}

// NewCollector creates a new Collector with the start time set to now.
func NewCollector(filterStages int) *Collector {
	return &Collector{
		StartTime:    time.Now(),
		FilterStages: filterStages,
	}
}

// Record updates counters for a single processed line.
func (c *Collector) Record(line string, matched bool) {
	c.LinesRead++
	c.BytesRead += int64(len(line)) + 1 // +1 for newline
	if matched {
		c.LinesMatched++
	} else {
		c.LinesDropped++
	}
}

// Elapsed returns the duration since the collector was created.
func (c *Collector) Elapsed() time.Duration {
	return time.Since(c.StartTime)
}

// MatchRate returns the fraction of lines that matched, in [0, 1].
func (c *Collector) MatchRate() float64 {
	if c.LinesRead == 0 {
		return 0
	}
	return float64(c.LinesMatched) / float64(c.LinesRead)
}

// Summary returns a human-readable summary string.
func (c *Collector) Summary() string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "lines read: %d\n", c.LinesRead)
	fmt.Fprintf(&sb, "lines matched: %d\n", c.LinesMatched)
	fmt.Fprintf(&sb, "lines dropped: %d\n", c.LinesDropped)
	fmt.Fprintf(&sb, "bytes read: %d\n", c.BytesRead)
	fmt.Fprintf(&sb, "filter stages: %d\n", c.FilterStages)
	fmt.Fprintf(&sb, "match rate: %.1f%%\n", c.MatchRate()*100)
	fmt.Fprintf(&sb, "elapsed: %s", c.Elapsed().Round(time.Millisecond))
	return sb.String()
}
