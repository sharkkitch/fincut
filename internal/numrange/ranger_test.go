package numrange

import (
	"testing"
)

func fptr(f float64) *float64 { return &f }

func TestNew_EmptyPattern(t *testing.T) {
	_, err := New(Options{})
	if err == nil {
		t.Fatal("expected error for empty pattern")
	}
}

func TestNew_InvalidPattern(t *testing.T) {
	_, err := New(Options{Pattern: "[invalid"})
	if err == nil {
		t.Fatal("expected error for invalid regex")
	}
}

func TestNew_NoCaptureGroup(t *testing.T) {
	_, err := New(Options{Pattern: `\d+`})
	if err == nil {
		t.Fatal("expected error when no capture group")
	}
}

func TestNew_MinExceedsMax(t *testing.T) {
	_, err := New(Options{Pattern: `(\d+)`, Min: fptr(10), Max: fptr(5)})
	if err == nil {
		t.Fatal("expected error when min > max")
	}
}

func TestNew_ValidOptions(t *testing.T) {
	_, err := New(Options{Pattern: `val=(\d+)`, Min: fptr(0), Max: fptr(100)})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRanger_Apply_WithinRange(t *testing.T) {
	r, _ := New(Options{Pattern: `score=(\d+)`, Min: fptr(50), Max: fptr(80)})
	lines := []string{"score=45", "score=60", "score=80", "score=99", "no-match"}
	got := r.Apply(lines)
	if len(got) != 2 {
		t.Fatalf("expected 2 lines, got %d: %v", len(got), got)
	}
}

func TestRanger_Apply_NoUpperBound(t *testing.T) {
	r, _ := New(Options{Pattern: `n=(\d+)`, Min: fptr(10)})
	lines := []string{"n=5", "n=10", "n=200"}
	got := r.Apply(lines)
	if len(got) != 2 {
		t.Fatalf("expected 2, got %d", len(got))
	}
}

func TestRanger_Apply_Invert(t *testing.T) {
	r, _ := New(Options{Pattern: `v=(-?[\d.]+)`, Min: fptr(0), Max: fptr(100), Invert: true})
	lines := []string{"v=50", "v=-1", "v=200"}
	got := r.Apply(lines)
	if len(got) != 2 {
		t.Fatalf("expected 2 inverted lines, got %d: %v", len(got), got)
	}
}

func TestRanger_Apply_NoMatch_Skipped(t *testing.T) {
	r, _ := New(Options{Pattern: `val=(\d+)`, Min: fptr(0)})
	lines := []string{"nothing here", "val=5"}
	got := r.Apply(lines)
	if len(got) != 1 {
		t.Fatalf("expected 1, got %d", len(got))
	}
}

func TestFormatSummary(t *testing.T) {
	s := FormatSummary(100, 42)
	if s == "" {
		t.Fatal("expected non-empty summary")
	}
}
