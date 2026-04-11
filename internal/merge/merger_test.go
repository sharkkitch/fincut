package merge

import (
	"bytes"
	"strings"
	"testing"
)

func TestNewMerger_SortByTimeWithoutLayout(t *testing.T) {
	_, err := NewMerger(Options{SortByTime: true})
	if err == nil {
		t.Fatal("expected error when SortByTime is true and TimestampLayout is empty")
	}
}

func TestNewMerger_ValidOptions(t *testing.T) {
	m, err := NewMerger(Options{LabelSources: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m == nil {
		t.Fatal("expected non-nil Merger")
	}
}

func TestMerger_Merge_NoLabel(t *testing.T) {
	m, _ := NewMerger(Options{})
	sources := map[string][]string{
		"a.log": {"line one", "line two"},
	}
	var buf bytes.Buffer
	if err := m.Merge(sources, &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := buf.String()
	if !strings.Contains(got, "line one") || !strings.Contains(got, "line two") {
		t.Errorf("expected both lines in output, got: %q", got)
	}
}

func TestMerger_Merge_WithLabel(t *testing.T) {
	m, _ := NewMerger(Options{LabelSources: true})
	sources := map[string][]string{
		"svc.log": {"hello world"},
	}
	var buf bytes.Buffer
	_ = m.Merge(sources, &buf)
	if !strings.Contains(buf.String(), "[svc.log]") {
		t.Errorf("expected source label in output, got: %q", buf.String())
	}
}

func TestMerger_Merge_SortByTime(t *testing.T) {
	layout := "2006-01-02T15:04:05"
	m, _ := NewMerger(Options{SortByTime: true, TimestampLayout: layout})
	sources := map[string][]string{
		"b.log": {"2024-06-01T12:00:05 second"},
		"a.log": {"2024-06-01T12:00:01 first"},
	}
	var buf bytes.Buffer
	_ = m.Merge(sources, &buf)
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(lines))
	}
	if !strings.Contains(lines[0], "first") {
		t.Errorf("expected 'first' to appear before 'second', got: %v", lines)
	}
}

func TestMerger_Merge_EmptySources(t *testing.T) {
	m, _ := NewMerger(Options{})
	var buf bytes.Buffer
	if err := m.Merge(map[string][]string{}, &buf); err != nil {
		t.Fatalf("unexpected error on empty sources: %v", err)
	}
	if buf.Len() != 0 {
		t.Errorf("expected empty output, got: %q", buf.String())
	}
}
