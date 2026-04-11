package segment_test

import (
	"strings"
	"testing"

	"github.com/user/fincut/internal/segment"
)

func TestSegmenter_LargeInput_EvenSplit(t *testing.T) {
	const total = 100
	const window = 10
	lines := make([]string, total)
	for i := range lines {
		lines[i] = "log line"
	}
	seg, err := segment.NewSegmenter(segment.Options{WindowSize: window})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	segs, err := seg.Apply(lines)
	if err != nil {
		t.Fatalf("apply error: %v", err)
	}
	if len(segs) != total/window {
		t.Errorf("expected %d segments, got %d", total/window, len(segs))
	}
	for i, s := range segs {
		if len(s.Lines) != window {
			t.Errorf("segment %d: expected %d lines, got %d", i, window, len(s.Lines))
		}
	}
}

func TestSegmenter_LabelsAreUnique(t *testing.T) {
	lines := []string{"a", "b", "c", "d", "e", "f"}
	seg, _ := segment.NewSegmenter(segment.Options{WindowSize: 2})
	segs, _ := seg.Apply(lines)
	seen := map[string]bool{}
	for _, s := range segs {
		if seen[s.Label] {
			t.Errorf("duplicate label: %s", s.Label)
		}
		seen[s.Label] = true
	}
}

func TestSegmenter_ContentPreserved(t *testing.T) {
	lines := []string{"alpha", "beta", "gamma", "delta"}
	seg, _ := segment.NewSegmenter(segment.Options{WindowSize: 2})
	segs, _ := seg.Apply(lines)
	var got []string
	for _, s := range segs {
		got = append(got, s.Lines...)
	}
	if strings.Join(got, ",") != strings.Join(lines, ",") {
		t.Errorf("content mismatch: got %v", got)
	}
}
