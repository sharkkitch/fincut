// Package segment provides log segmentation by time window or record count.
package segment

import (
	"errors"
	"time"
)

// Segment represents a named slice of log lines.
type Segment struct {
	Label string
	Lines []string
}

// Options configures how segmentation is performed.
type Options struct {
	// WindowSize groups lines into fixed-size buckets.
	WindowSize int
	// TimestampLayout is the Go time layout used to parse timestamps.
	// When set, segmentation is based on time windows of Duration.
	TimestampLayout string
	Duration        time.Duration
	// ExtractTimestamp is called per line to extract a timestamp substring.
	ExtractTimestamp func(line string) string
}

// Segmenter splits a stream of lines into labelled segments.
type Segmenter struct {
	opts Options
}

// NewSegmenter validates opts and returns a ready Segmenter.
func NewSegmenter(opts Options) (*Segmenter, error) {
	if opts.WindowSize < 0 {
		return nil, errors.New("segment: WindowSize must be >= 0")
	}
	if opts.TimestampLayout != "" && opts.Duration <= 0 {
		return nil, errors.New("segment: Duration must be > 0 when TimestampLayout is set")
	}
	if opts.TimestampLayout != "" && opts.ExtractTimestamp == nil {
		return nil, errors.New("segment: ExtractTimestamp func required when TimestampLayout is set")
	}
	return &Segmenter{opts: opts}, nil
}

// Apply partitions lines into segments according to the configured strategy.
func (s *Segmenter) Apply(lines []string) ([]Segment, error) {
	if s.opts.TimestampLayout != "" {
		return s.byTime(lines)
	}
	size := s.opts.WindowSize
	if size == 0 {
		size = len(lines)
		if size == 0 {
			return nil, nil
		}
	}
	return s.byCount(lines, size), nil
}

func (s *Segmenter) byCount(lines []string, size int) []Segment {
	var segs []Segment
	for i := 0; i < len(lines); i += size {
		end := i + size
		if end > len(lines) {
			end = len(lines)
		}
		segs = append(segs, Segment{
			Label: formatLabel(i/size + 1),
			Lines: lines[i:end],
		})
	}
	return segs
}

func (s *Segmenter) byTime(lines []string) ([]Segment, error) {
	var (
		segs    []Segment
		current *Segment
		windowStart time.Time
	)
	for _, line := range lines {
		raw := s.opts.ExtractTimestamp(line)
		t, err := time.Parse(s.opts.TimestampLayout, raw)
		if err != nil {
			if current != nil {
				current.Lines = append(current.Lines, line)
			}
			continue
		}
		if current == nil || t.Sub(windowStart) >= s.opts.Duration {
			windowStart = t
			segs = append(segs, Segment{Label: t.Format(s.opts.TimestampLayout)})
			current = &segs[len(segs)-1]
		}
		current.Lines = append(current.Lines, line)
	}
	return segs, nil
}

func formatLabel(n int) string {
	return "segment-" + itoa(n)
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	buf := make([]byte, 0, 10)
	for n > 0 {
		buf = append([]byte{byte('0' + n%10)}, buf...)
		n /= 10
	}
	return string(buf)
}
