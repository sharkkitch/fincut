package linenum

import (
	"testing"
)

func TestNew_NoRanges(t *testing.T) {
	_, err := New(Options{})
	if err == nil {
		t.Fatal("expected error for empty ranges")
	}
}

func TestNew_InvalidStart(t *testing.T) {
	_, err := New(Options{Ranges: [][2]int{{0, 5}}})
	if err == nil {
		t.Fatal("expected error for start < 1")
	}
}

func TestNew_EndBeforeStart(t *testing.T) {
	_, err := New(Options{Ranges: [][2]int{{5, 3}}})
	if err == nil {
		t.Fatal("expected error for end < start")
	}
}

func TestNew_ValidOptions(t *testing.T) {
	_, err := New(Options{Ranges: [][2]int{{1, 5}}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestLinenupper_Apply_SingleRange(t *testing.T) {
	l, _ := New(Options{Ranges: [][2]int{{2, 4}}})
	input := []string{"a", "b", "c", "d", "e"}
	got := l.Apply(input)
	want := []string{"b", "c", "d"}
	if len(got) != len(want) {
		t.Fatalf("want %v, got %v", want, got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("index %d: want %q, got %q", i, want[i], got[i])
		}
	}
}

func TestLinenupper_Apply_MultipleRanges(t *testing.T) {
	l, _ := New(Options{Ranges: [][2]int{{1, 1}, {4, 5}}})
	input := []string{"a", "b", "c", "d", "e"}
	got := l.Apply(input)
	want := []string{"a", "d", "e"}
	if len(got) != len(want) {
		t.Fatalf("want %v, got %v", want, got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("index %d: want %q, got %q", i, want[i], got[i])
		}
	}
}

func TestLinenupper_Apply_EOFRange(t *testing.T) {
	l, _ := New(Options{Ranges: [][2]int{{3, 0}}})
	input := []string{"a", "b", "c", "d", "e"}
	got := l.Apply(input)
	want := []string{"c", "d", "e"}
	if len(got) != len(want) {
		t.Fatalf("want %v, got %v", want, got)
	}
}

func TestLinenupper_Apply_EmptyInput(t *testing.T) {
	l, _ := New(Options{Ranges: [][2]int{{1, 5}}})
	got := l.Apply([]string{})
	if len(got) != 0 {
		t.Fatalf("expected empty output, got %v", got)
	}
}

func TestFormatSummary_NoRanges(t *testing.T) {
	s := FormatSummary(Options{})
	if s != "line ranges: none" {
		t.Errorf("unexpected summary: %q", s)
	}
}

func TestFormatSummary_WithRanges(t *testing.T) {
	s := FormatSummary(Options{Ranges: [][2]int{{1, 3}, {7, 0}}})
	if s == "" {
		t.Error("expected non-empty summary")
	}
}
