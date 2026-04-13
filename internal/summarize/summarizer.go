package summarize

import (
	"fmt"
	"strings"
)

// Options configures the Summarizer.
type Options struct {
	// TopN limits the number of top patterns reported. Zero means no limit.
	TopN int
	// MinCount excludes patterns seen fewer than MinCount times.
	MinCount int
	// CaseSensitive controls whether counting is case-sensitive.
	CaseSensitive bool
}

// entry tracks an observed line and its frequency.
type entry struct {
	line  string
	count int
}

// Summarizer counts distinct lines and reports frequency summaries.
type Summarizer struct {
	opts    Options
	centry
	origKey map[string]string // normalised key -> first-seen original
	total   int
}

// New constructs a Summarizer. Returns an error if options are invalid.
func New(opts Options) (*Summarizer, error) {
	if err := validateOptions(opts); err != nil {
		return nil, err
	}
	return &Summarizer{
		opts:    opts,
		counts:  make(map[string]*entry),
		origKey: make(map[string]string),
	}, nil
}

// Add records a single line.
func (s *Summarizer) Add(line string) {
	s.total++
	key := line
	if !s.opts.CaseSensitive {
		key = strings.ToLower(line)
	}
	if e, ok := s.counts[key]; ok {
		e.count++
	} else {
		s.counts[key] = &entry{line: line, count: 1}
	}
}

// Results returns entries sorted by descending count, filtered by MinCount,
// and limited to TopN (if set).
func (s *Summarizer) Results() []Result {
	var out []Result
	for _, e := range s.counts {
		if e.count >= minCount(s.opts.MinCount) {
			out = append(out, Result{Line: e.line, Count: e.count})
		}
	}
	sortResults(out)
	if s.opts.TopN > 0 && len(out) > s.opts.TopN {
		out = out[:s.opts.TopN]
	}
	return out
}

// Total returns the total number of lines added.
func (s *Summarizer) Total() int { return s.total }

func minCount(n int) int {
	if n <= 0 {
		return 1
	}
	return n
}

// FormatSummary returns a human-readable summary string.
func FormatSummary(results []Result, total int) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "total lines: %d, distinct: %d\n", total, len(results))
	for _, r := range results {
		pct := 0.0
		if total > 0 {
			pct = float64(r.Count) / float64(total) * 100
		}
		fmt.Fprintf(&sb, "  [%5d | %5.1f%%] %s\n", r.Count, pct, r.Line)
	}
	return sb.String()
}
