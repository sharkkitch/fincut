package segment

import (
	"testing"
	"time"
)

func TestNewSegmenter_InvalidOptions(t *testing.T) {
	cases := []struct {
		name string
		opts Options
	}{
		{"negative window", Options{WindowSize: -1}},
		{"layout without duration", Options{TimestampLayout: time.RFC3339, ExtractTimestamp: func(s string) string { return s }}},
		{"layout without extractor", Options{TimestampLayout: time.RFC3339, Duration: time.Minute}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := NewSegmenter(tc.opts)
			if err == nil {
				t.Fatal("expected error, got nil")
			}
		})
	}
}

func TestNewSegmenter_ValidOptions(t *testing.T) {
	_, err := NewSegmenter(Options{WindowSize: 10})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestSegmenter_Apply_ByCount(t *testing.T) {
	seg, _ := NewSegmenter(Options{WindowSize: 3})
	lines := []string{"a", "b", "c", "d", "e"}
	segs, err := seg.Apply(lines)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(segs) != 2 {
		t.Fatalf("expected 2 segments, got %d", len(segs))
	}
	if len(segs[0].Lines) != 3 {
		t.Errorf("segment 0: expected 3 lines, got %d", len(segs[0].Lines))
	}
	if len(segs[1].Lines) != 2 {
		t.Errorf("segment 1: expected 2 lines, got %d", len(segs[1].Lines))
	}
}

func TestSegmenter_Apply_EmptyInput(t *testing.T) {
	seg, _ := NewSegmenter(Options{WindowSize: 5})
	segs, err := seg.Apply([]string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(segs) != 0 {
		t.Errorf("expected 0 segments, got %d", len(segs))
	}
}

func TestSegmenter_Apply_SingleSegment(t *testing.T) {
	seg, _ := NewSegmenter(Options{WindowSize: 0})
	lines := []string{"x", "y", "z"}
	segs, err := seg.Apply(lines)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(segs) != 1 {
		t.Fatalf("expected 1 segment, got %d", len(segs))
	}
	if len(segs[0].Lines) != 3 {
		t.Errorf("expected 3 lines, got %d", len(segs[0].Lines))
	}
}

func TestSegmenter_Apply_ByTime(t *testing.T) {
	layout := "2006-01-02T15:04:05"
	extract := func(line string) string {
		if len(line) >= len(layout) {
			return line[:len(layout)]
		}
		return ""
	}
	seg, _ := NewSegmenter(Options{
		TimestampLayout:  layout,
		Duration:         time.Minute,
		ExtractTimestamp: extract,
	})
	lines := []string{
		"2024-01-01T10:00:00 line1",
		"2024-01-01T10:00:30 line2",
		"2024-01-01T10:01:05 line3",
	}
	segs, err := seg.Apply(lines)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(segs) != 2 {
		t.Fatalf("expected 2 time segments, got %d", len(segs))
	}
	if len(segs[0].Lines) != 2 {
		t.Errorf("expected 2 lines in first segment, got %d", len(segs[0].Lines))
	}
}
