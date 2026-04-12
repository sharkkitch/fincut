// Package freq provides line frequency analysis for structured log files.
// It counts how often each unique line (or field) appears and returns
// results ranked by occurrence count.
package freq

import (
	"fmt"
	"sort"
	"strings"
)

// Options configures the frequency counter.
type Options struct {
	// TopN limits results to the N most frequent entries. 0 means no limit.
	TopN int
	// CaseInsensitive normalises lines to lowercase before counting.
	CaseInsensitive bool
	// Field selects a whitespace-delimited field index (1-based). 0 means whole line.
	Field int
}

// Entry holds a single frequency result.
type Entry struct {
	Value string
	Count int
}

// Counter accumulates line frequencies.
type Counter struct {
	opts   Options
	counts map[string]int
}

// New creates a Counter from opts, returning an error for invalid options.
func New(opts Options) (*Counter, error) {
	if opts.TopN < 0 {
		return nil, fmt.Errorf("freq: TopN must be >= 0, got %d", opts.TopN)
	}
	if opts.Field < 0 {
		return nil, fmt.Errorf("freq: Field must be >= 0, got %d", opts.Field)
	}
	return &Counter{opts: opts, counts: make(map[string]int)}, nil
}

// Add records each line in lines into the frequency table.
func (c *Counter) Add(lines []string) {
	for _, l := range lines {
		key := c.key(l)
		c.counts[key]++
	}
}

// Results returns entries sorted descending by count, capped to TopN if set.
func (c *Counter) Results() []Entry {
	entries := make([]Entry, 0, len(c.counts))
	for v, n := range c.counts {
		entries = append(entries, Entry{Value: v, Count: n})
	}
	sort.Slice(entries, func(i, j int) bool {
		if entries[i].Count != entries[j].Count {
			return entries[i].Count > entries[j].Count
		}
		return entries[i].Value < entries[j].Value
	})
	if c.opts.TopN > 0 && len(entries) > c.opts.TopN {
		entries = entries[:c.opts.TopN]
	}
	return entries
}

// Total returns the total number of lines processed.
func (c *Counter) Total() int {
	n := 0
	for _, v := range c.counts {
		n += v
	}
	return n
}

func (c *Counter) key(line string) string {
	s := line
	if c.opts.Field > 0 {
		parts := strings.Fields(s)
		if c.opts.Field <= len(parts) {
			s = parts[c.opts.Field-1]
		} else {
			s = ""
		}
	}
	if c.opts.CaseInsensitive {
		s = strings.ToLower(s)
	}
	return s
}
