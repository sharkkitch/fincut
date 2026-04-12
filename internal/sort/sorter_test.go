package sort

import (
	"testing"
)

func TestNew_NegativeField(t *testing.T) {
	_, err := New(Options{Field: -1})
	if err == nil {
		t.Fatal("expected error for negative field")
	}
}

func TestNew_FieldWithoutDelimiter(t *testing.T) {
	_, err := New(Options{Field: 2})
	if err == nil {
		t.Fatal("expected error: field > 0 requires delimiter")
	}
}

func TestNew_ValidOptions(t *testing.T) {
	s, err := New(Options{Reverse: true, Unique: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s == nil {
		t.Fatal("expected non-nil sorter")
	}
}

func TestSorter_Apply_Ascending(t *testing.T) {
	s, _ := New(Options{})
	in := []string{"banana", "apple", "cherry"}
	out := s.Apply(in)
	want := []string{"apple", "banana", "cherry"}
	for i, v := range want {
		if out[i] != v {
			t.Errorf("pos %d: got %q, want %q", i, out[i], v)
		}
	}
}

func TestSorter_Apply_Descending(t *testing.T) {
	s, _ := New(Options{Reverse: true})
	in := []string{"banana", "apple", "cherry"}
	out := s.Apply(in)
	if out[0] != "cherry" {
		t.Errorf("expected cherry first, got %q", out[0])
	}
}

func TestSorter_Apply_Unique(t *testing.T) {
	s, _ := New(Options{Unique: true})
	in := []string{"b", "a", "b", "c", "a"}
	out := s.Apply(in)
	if len(out) != 3 {
		t.Errorf("expected 3 unique lines, got %d", len(out))
	}
}

func TestSorter_Apply_ByField(t *testing.T) {
	s, _ := New(Options{Delimiter: "|", Field: 2})
	in := []string{"row|banana|1", "row|apple|2", "row|cherry|3"}
	out := s.Apply(in)
	if out[0] != "row|apple|2" {
		t.Errorf("expected apple row first, got %q", out[0])
	}
}

func TestSorter_Apply_EmptyInput(t *testing.T) {
	s, _ := New(Options{})
	out := s.Apply([]string{})
	if len(out) != 0 {
		t.Errorf("expected empty output, got %d lines", len(out))
	}
}

func TestSorter_Apply_FieldOutOfRange(t *testing.T) {
	s, _ := New(Options{Delimiter: "|", Field: 5})
	in := []string{"a|b", "c|d"}
	out := s.Apply(in)
	// all keys are empty string — order stable, length preserved
	if len(out) != 2 {
		t.Errorf("expected 2 lines, got %d", len(out))
	}
}
