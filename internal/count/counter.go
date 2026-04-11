// Package count provides line counting and frequency analysis for log streams.
package count

import (
	"fmt"
	"sort"
	"strings"
)

// Entry holds a line value and the number of times it was observed.
type Entry struct {
	Line  string
	Count int
}

// Options configures the Counter behaviour.
type Options struct {
	// TopN limits output to the N most frequent lines. Zero means unlimited.
	TopN int
	// CaseInsensitive folds lines to lowercase before counting.
	CaseInsensitive bool
}

// Counter tallies occurrences of each unique line.
type Counter struct {
	opts   Options
	counts map[string]int
	// canonical stores the first-seen casing for a key.
	canonical map[string]string
}

// New returns a Counter configured with opts, or an error if opts are invalid.
func New(opts Options) (*Counter, error) {
	if opts.TopN < 0 {
		return nil, fmt.Errorf("count: TopN must be >= 0, got %d", opts.TopN)
	}
	return &Counter{
		opts:      opts,
		counts:    make(map[string]int),
		canonical: make(map[string]string),
	}, nil
}

// Add records a single line occurrence.
func (c *Counter) Add(line string) {
	key := line
	if c.opts.CaseInsensitive {
		key = strings.ToLower(line)
	}
	if _, seen := c.canonical[key]; !seen {
		c.canonical[key] = line
	}
	c.counts[key]++
}

// Results returns entries sorted by descending frequency.
// If TopN > 0, only the top N entries are returned.
func (c *Counter) Results() []Entry {
	entries := make([]Entry, 0, len(c.counts))
	for key, n := range c.counts {
		entries = append(entries, Entry{Line: c.canonical[key], Count: n})
	}
	sort.Slice(entries, func(i, j int) bool {
		if entries[i].Count != entries[j].Count {
			return entries[i].Count > entries[j].Count
		}
		return entries[i].Line < entries[j].Line
	})
	if c.opts.TopN > 0 && len(entries) > c.opts.TopN {
		entries = entries[:c.opts.TopN]
	}
	return entries
}

// Total returns the total number of lines added.
func (c *Counter) Total() int {
	total := 0
	for _, n := range c.counts {
		total += n
	}
	return total
}
