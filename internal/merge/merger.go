package merge

import (
	"fmt"
	"io"
	"sort"
	"time"
)

// Line represents a single log line with its source file and optional timestamp.
type Line struct {
	Text      string
	Source    string
	Timestamp time.Time
	Index     int
}

// Options configures the merge behaviour.
type Options struct {
	// SortByTime sorts merged output by parsed timestamp; requires TimestampLayout.
	SortByTime bool
	// TimestampLayout is the Go time layout used to parse timestamps from lines.
	TimestampLayout string
	// LabelSources prefixes each output line with its source filename.
	LabelSources bool
}

// Merger merges multiple slices of log lines into a single ordered output.
type Merger struct {
	opts Options
}

// NewMerger creates a Merger from the given Options.
func NewMerger(opts Options) (*Merger, error) {
	if opts.SortByTime && opts.TimestampLayout == "" {
		return nil, fmt.Errorf("merge: TimestampLayout required when SortByTime is true")
	}
	return &Merger{opts: opts}, nil
}

// Merge combines lines from all sources and writes them to w.
func (m *Merger) Merge(sources map[string][]string, w io.Writer) error {
	var lines []Line

	for src, texts := range sources {
		for i, text := range texts {
			l := Line{Text: text, Source: src, Index: i}
			if m.opts.SortByTime && m.opts.TimestampLayout != "" {
				ts, err := parseTimestamp(text, m.opts.TimestampLayout)
				if err == nil {
					l.Timestamp = ts
				}
			}
			lines = append(lines, l)
		}
	}

	if m.opts.SortByTime {
		sort.SliceStable(lines, func(i, j int) bool {
			return lines[i].Timestamp.Before(lines[j].Timestamp)
		})
	}

	for _, l := range lines {
		out := l.Text
		if m.opts.LabelSources {
			out = fmt.Sprintf("[%s] %s", l.Source, l.Text)
		}
		if _, err := fmt.Fprintln(w, out); err != nil {
			return fmt.Errorf("merge: write error: %w", err)
		}
	}
	return nil
}

// parseTimestamp attempts to find and parse a timestamp within text.
func parseTimestamp(text, layout string) (time.Time, error) {
	if len(text) >= len(layout) {
		return time.Parse(layout, text[:len(layout)])
	}
	return time.Time{}, fmt.Errorf("merge: text too short for layout")
}
