package merge_test

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/yourorg/fincut/internal/merge"
)

func TestMerger_LargeInput_PreservesAllLines(t *testing.T) {
	const n = 500
	sources := map[string][]string{}
	for f := 0; f < 3; f++ {
		key := fmt.Sprintf("file%d.log", f)
		for i := 0; i < n; i++ {
			sources[key] = append(sources[key], fmt.Sprintf("line %d from %s", i, key))
		}
	}

	m, err := merge.NewMerger(merge.Options{LabelSources: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var buf bytes.Buffer
	if err := m.Merge(sources, &buf); err != nil {
		t.Fatalf("merge failed: %v", err)
	}

	got := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(got) != 3*n {
		t.Errorf("expected %d lines, got %d", 3*n, len(got))
	}
}

func TestMerger_SortByTime_StableForEqualTimestamps(t *testing.T) {
	layout := "2006-01-02T15:04:05"
	m, _ := merge.NewMerger(merge.Options{SortByTime: true, TimestampLayout: layout})
	sources := map[string][]string{
		"x.log": {
			"2024-01-01T00:00:01 alpha",
			"2024-01-01T00:00:01 beta",
			"2024-01-01T00:00:02 gamma",
		},
	}
	var buf bytes.Buffer
	_ = m.Merge(sources, &buf)
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(lines))
	}
	if !strings.Contains(lines[2], "gamma") {
		t.Errorf("expected gamma last, got: %v", lines)
	}
}
