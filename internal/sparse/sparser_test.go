package sparse

import (
	"testing"
)

func TestNew_EveryTooSmall(t *testing.T) {
	_, err := New(Options{Every: 0, Offset: 0})
	if err == nil {
		t.Fatal("expected error for Every=0")
	}
}

func TestNew_NegativeOffset(t *testing.T) {
	_, err := New(Options{Every: 2, Offset: -1})
	if err == nil {
		t.Fatal("expected error for negative Offset")
	}
}

func TestNew_OffsetExceedsEvery(t *testing.T) {
	_, err := New(Options{Every: 2, Offset: 2})
	if err == nil {
		t.Fatal("expected error when Offset >= Every")
	}
}

func TestNew_ValidOptions(t *testing.T) {
	s, err := New(Options{Every: 3, Offset: 1})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s == nil {
		t.Fatal("expected non-nil Sparser")
	}
}

func TestSparser_Apply_EmptyInput(t *testing.T) {
	s, _ := New(Options{Every: 2, Offset: 0})
	out := s.Apply([]string{})
	if len(out) != 0 {
		t.Fatalf("expected empty output, got %d lines", len(out))
	}
}

func TestSparser_Apply_EverySecond(t *testing.T) {
	lines := []string{"a", "b", "c", "d", "e", "f"}
	s, _ := New(Options{Every: 2, Offset: 0})
	out := s.Apply(lines)
	want := []string{"a", "c", "e"}
	if len(out) != len(want) {
		t.Fatalf("expected %d lines, got %d", len(want), len(out))
	}
	for i, w := range want {
		if out[i] != w {
			t.Errorf("line %d: want %q got %q", i, w, out[i])
		}
	}
}

func TestSparser_Apply_WithOffset(t *testing.T) {
	lines := []string{"a", "b", "c", "d", "e", "f"}
	s, _ := New(Options{Every: 2, Offset: 1})
	out := s.Apply(lines)
	want := []string{"b", "d", "f"}
	if len(out) != len(want) {
		t.Fatalf("expected %d lines, got %d", len(want), len(out))
	}
	for i, w := range want {
		if out[i] != w {
			t.Errorf("line %d: want %q got %q", i, w, out[i])
		}
	}
}

func TestSparser_Apply_Every1_KeepsAll(t *testing.T) {
	lines := []string{"x", "y", "z"}
	s, _ := New(Options{Every: 1, Offset: 0})
	out := s.Apply(lines)
	if len(out) != len(lines) {
		t.Fatalf("expected all %d lines, got %d", len(lines), len(out))
	}
}

func TestFormatSummary(t *testing.T) {
	summary := FormatSummary(100, 34, Options{Every: 3, Offset: 1})
	if summary == "" {
		t.Fatal("expected non-empty summary")
	}
}
