package head

import (
	"strings"
	"testing"
)

func TestNew_NegativeMaxLines(t *testing.T) {
	_, err := New(Options{MaxLines: -1})
	if err == nil {
		t.Fatal("expected error for negative MaxLines")
	}
}

func TestNew_NegativeMaxBytes(t *testing.T) {
	_, err := New(Options{MaxBytes: -1})
	if err == nil {
		t.Fatal("expected error for negative MaxBytes")
	}
}

func TestNew_NeitherSet(t *testing.T) {
	_, err := New(Options{})
	if err == nil {
		t.Fatal("expected error when neither MaxLines nor MaxBytes is set")
	}
}

func TestNew_ValidOptions(t *testing.T) {
	h, err := New(Options{MaxLines: 10})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if h == nil {
		t.Fatal("expected non-nil Header")
	}
}

func TestHeader_Apply_MaxLines(t *testing.T) {
	h, _ := New(Options{MaxLines: 3})
	input := []string{"a", "b", "c", "d", "e"}
	got := h.Apply(input)
	if len(got) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(got))
	}
	if strings.Join(got, ",") != "a,b,c" {
		t.Fatalf("unexpected result: %v", got)
	}
}

func TestHeader_Apply_FewerLinesThanMax(t *testing.T) {
	h, _ := New(Options{MaxLines: 100})
	input := []string{"x", "y"}
	got := h.Apply(input)
	if len(got) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(got))
	}
}

func TestHeader_Apply_MaxBytes(t *testing.T) {
	// Each line is 3 bytes + 1 newline = 4 bytes. Limit 9 allows 2 lines (8 bytes).
	h, _ := New(Options{MaxBytes: 9})
	input := []string{"abc", "def", "ghi"}
	got := h.Apply(input)
	if len(got) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(got))
	}
}

func TestHeader_Apply_BothLimits_LineWins(t *testing.T) {
	// MaxLines=1 is stricter than MaxBytes=1000
	h, _ := New(Options{MaxLines: 1, MaxBytes: 1000})
	input := []string{"line1", "line2", "line3"}
	got := h.Apply(input)
	if len(got) != 1 {
		t.Fatalf("expected 1 line, got %d", len(got))
	}
}

func TestHeader_Apply_EmptyInput(t *testing.T) {
	h, _ := New(Options{MaxLines: 5})
	got := h.Apply([]string{})
	if len(got) != 0 {
		t.Fatalf("expected empty result, got %v", got)
	}
}
