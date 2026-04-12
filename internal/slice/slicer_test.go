package slice

import (
	"strings"
	"testing"
)

func makeLines(n int) []string {
	lines := make([]string, n)
	for i := range lines {
		lines[i] = strings.Repeat("x", i+1) // unique content per line
	}
	return lines
}

func TestNew_InvalidEnd(t *testing.T) {
	_, err := New(Options{Start: 5, End: 3})
	if err == nil {
		t.Fatal("expected error for End < Start")
	}
}

func TestNew_NegativeStep(t *testing.T) {
	_, err := New(Options{Step: -1})
	if err == nil {
		t.Fatal("expected error for negative Step")
	}
}

func TestNew_ValidOptions(t *testing.T) {
	s, err := New(Options{Start: 2, End: 8, Step: 2})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s == nil {
		t.Fatal("expected non-nil Slicer")
	}
}

func TestSlicer_Apply_AllLines(t *testing.T) {
	s, _ := New(Options{})
	lines := makeLines(5)
	got := s.Apply(lines)
	if len(got) != 5 {
		t.Fatalf("expected 5 lines, got %d", len(got))
	}
}

func TestSlicer_Apply_StartOffset(t *testing.T) {
	s, _ := New(Options{Start: 3})
	lines := makeLines(5)
	got := s.Apply(lines)
	if len(got) != 3 {
		t.Fatalf("expected 3 lines (3-5), got %d", len(got))
	}
}

func TestSlicer_Apply_EndOffset(t *testing.T) {
	s, _ := New(Options{End: 3})
	lines := makeLines(5)
	got := s.Apply(lines)
	if len(got) != 3 {
		t.Fatalf("expected 3 lines (1-3), got %d", len(got))
	}
}

func TestSlicer_Apply_Range(t *testing.T) {
	s, _ := New(Options{Start: 2, End: 4})
	lines := makeLines(6)
	got := s.Apply(lines)
	if len(got) != 3 {
		t.Fatalf("expected 3 lines (2-4), got %d", len(got))
	}
	if got[0] != lines[1] {
		t.Errorf("first result should be line 2")
	}
}

func TestSlicer_Apply_Step(t *testing.T) {
	s, _ := New(Options{Start: 1, End: 10, Step: 3})
	lines := makeLines(10)
	got := s.Apply(lines)
	// positions 1,4,7,10 => 4 lines
	if len(got) != 4 {
		t.Fatalf("expected 4 lines with step 3, got %d", len(got))
	}
}

func TestSlicer_Apply_EmptyInput(t *testing.T) {
	s, _ := New(Options{Start: 1, End: 5})
	got := s.Apply([]string{})
	if len(got) != 0 {
		t.Fatalf("expected 0 lines, got %d", len(got))
	}
}

func TestSlicer_Apply_StartBeyondInput(t *testing.T) {
	s, _ := New(Options{Start: 20})
	lines := makeLines(5)
	got := s.Apply(lines)
	if len(got) != 0 {
		t.Fatalf("expected 0 lines when Start > len, got %d", len(got))
	}
}
