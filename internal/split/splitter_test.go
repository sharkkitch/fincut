package split

import (
	"strings"
	"testing"
)

func TestNew_EmptyPattern(t *testing.T) {
	_, err := New(Options{})
	if err == nil {
		t.Fatal("expected error for empty pattern")
	}
}

func TestNew_InvalidPattern(t *testing.T) {
	_, err := New(Options{Pattern: `[invalid`})
	if err == nil {
		t.Fatal("expected error for invalid regex")
	}
}

func TestNew_ValidOptions(t *testing.T) {
	s, err := New(Options{Pattern: `^---`})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s == nil {
		t.Fatal("expected non-nil splitter")
	}
}

func TestSplitter_Apply_NoDelimiter(t *testing.T) {
	s, _ := New(Options{Pattern: `^---`})
	lines := []string{"alpha", "beta", "gamma"}
	groups := s.Apply(lines)
	if len(groups) != 1 {
		t.Fatalf("expected 1 group, got %d", len(groups))
	}
	if len(groups[0].Lines) != 3 {
		t.Fatalf("expected 3 lines in group, got %d", len(groups[0].Lines))
	}
}

func TestSplitter_Apply_SplitsOnPattern(t *testing.T) {
	s, _ := New(Options{Pattern: `^---`})
	lines := []string{"a", "b", "---", "c", "d"}
	groups := s.Apply(lines)
	if len(groups) != 2 {
		t.Fatalf("expected 2 groups, got %d", len(groups))
	}
	if len(groups[0].Lines) != 2 || len(groups[1].Lines) != 2 {
		t.Fatalf("unexpected line distribution: %v", groups)
	}
}

func TestSplitter_Apply_KeepDelimiter(t *testing.T) {
	s, _ := New(Options{Pattern: `^---`, KeepDelimiter: true})
	lines := []string{"a", "---", "b"}
	groups := s.Apply(lines)
	if len(groups) != 2 {
		t.Fatalf("expected 2 groups, got %d", len(groups))
	}
	if groups[1].Lines[0] != "---" {
		t.Fatalf("expected delimiter as first line of second group")
	}
}

func TestSplitter_Apply_EmptyInput(t *testing.T) {
	s, _ := New(Options{Pattern: `^---`})
	groups := s.Apply([]string{})
	if len(groups) != 1 {
		t.Fatalf("expected 1 empty group, got %d", len(groups))
	}
	if len(groups[0].Lines) != 0 {
		t.Fatal("expected empty group")
	}
}

func TestSplitter_Apply_CustomLabel(t *testing.T) {
	s, _ := New(Options{Pattern: `^===`, Label: "chunk"})
	lines := []string{"x", "===", "y"}
	groups := s.Apply(lines)
	for _, g := range groups {
		if !strings.HasPrefix(g.Label, "chunk-") {
			t.Fatalf("expected label prefix 'chunk-', got %q", g.Label)
		}
	}
}

func TestFormatSummary(t *testing.T) {
	groups := []Group{
		{Label: "group-1", Lines: []string{"a", "b"}},
		{Label: "group-2", Lines: []string{"c"}},
	}
	summary := FormatSummary(groups)
	if !strings.Contains(summary, "2 group") {
		t.Fatalf("unexpected summary: %q", summary)
	}
	if !strings.Contains(summary, "3 line") {
		t.Fatalf("unexpected summary: %q", summary)
	}
}
