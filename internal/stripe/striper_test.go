package stripe

import (
	"testing"
)

func TestNew_EveryTooSmall(t *testing.T) {
	_, err := New(Options{Every: 1, Offset: 0})
	if err == nil {
		t.Fatal("expected error for every=1")
	}
}

func TestNew_NegativeOffset(t *testing.T) {
	_, err := New(Options{Every: 3, Offset: -1})
	if err == nil {
		t.Fatal("expected error for negative offset")
	}
}

func TestNew_OffsetExceedsEvery(t *testing.T) {
	_, err := New(Options{Every: 3, Offset: 3})
	if err == nil {
		t.Fatal("expected error when offset >= every")
	}
}

func TestNew_ValidOptions(t *testing.T) {
	s, err := New(Options{Every: 4, Offset: 1})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s == nil {
		t.Fatal("expected non-nil Striper")
	}
}

func TestStriper_Apply_EmptyInput(t *testing.T) {
	s, _ := New(Options{Every: 2, Offset: 0})
	out := s.Apply(nil)
	if len(out) != 0 {
		t.Fatalf("expected empty output, got %d lines", len(out))
	}
}

func TestStriper_Apply_EveryTwo_OffsetZero(t *testing.T) {
	lines := []string{"a", "b", "c", "d", "e", "f"}
	s, _ := New(Options{Every: 2, Offset: 0})
	out := s.Apply(lines)
	want := []string{"a", "c", "e"}
	if len(out) != len(want) {
		t.Fatalf("expected %d lines, got %d", len(want), len(out))
	}
	for i, w := range want {
		if out[i] != w {
			t.Errorf("line %d: want %q, got %q", i, w, out[i])
		}
	}
}

func TestStriper_Apply_EveryTwo_OffsetOne(t *testing.T) {
	lines := []string{"a", "b", "c", "d", "e", "f"}
	s, _ := New(Options{Every: 2, Offset: 1})
	out := s.Apply(lines)
	want := []string{"b", "d", "f"}
	if len(out) != len(want) {
		t.Fatalf("expected %d lines, got %d", len(want), len(out))
	}
	for i, w := range want {
		if out[i] != w {
			t.Errorf("line %d: want %q, got %q", i, w, out[i])
		}
	}
}

func TestStriper_Apply_EveryThree(t *testing.T) {
	lines := []string{"0", "1", "2", "3", "4", "5", "6", "7", "8"}
	s, _ := New(Options{Every: 3, Offset: 2})
	out := s.Apply(lines)
	want := []string{"2", "5", "8"}
	if len(out) != len(want) {
		t.Fatalf("expected %d lines, got %d", len(want), len(out))
	}
	for i, w := range want {
		if out[i] != w {
			t.Errorf("line %d: want %q, got %q", i, w, out[i])
		}
	}
}

func TestFormatSummary(t *testing.T) {
	got := FormatSummary(9, 3, 3, 0)
	if got == "" {
		t.Fatal("expected non-empty summary")
	}
}
