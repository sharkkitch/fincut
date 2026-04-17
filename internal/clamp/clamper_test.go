package clamp

import (
	"testing"
)

func fptr(f float64) *float64 { return &f }

func TestNew_EmptyPattern(t *testing.T) {
	_, err := New(Options{Min: fptr(0), Max: fptr(100)})
	if err == nil {
		t.Fatal("expected error for empty pattern")
	}
}

func TestNew_NeitherSet(t *testing.T) {
	_, err := New(Options{Pattern: `val=(\d+)`})
	if err == nil {
		t.Fatal("expected error when neither Min nor Max set")
	}
}

func TestNew_MinExceedsMax(t *testing.T) {
	_, err := New(Options{Pattern: `val=(\d+)`, Min: fptr(100), Max: fptr(10)})
	if err == nil {
		t.Fatal("expected error when min > max")
	}
}

func TestNew_InvalidPattern(t *testing.T) {
	_, err := New(Options{Pattern: `[invalid`, Min: fptr(0)})
	if err == nil {
		t.Fatal("expected error for invalid regex")
	}
}

func TestNew_NoCaptureGroup(t *testing.T) {
	_, err := New(Options{Pattern: `val=\d+`, Min: fptr(0)})
	if err == nil {
		t.Fatal("expected error when pattern has no capture group")
	}
}

func TestNew_ValidOptions(t *testing.T) {
	_, err := New(Options{Pattern: `val=(\d+)`, Min: fptr(0), Max: fptr(100)})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestClamper_Apply_BasicRange(t *testing.T) {
	c, _ := New(Options{Pattern: `status=(\d+)`, Min: fptr(200), Max: fptr(299)})
	lines := []string{"status=200 OK", "status=404 Not Found", "status=201 Created", "no match here"}
	got := c.Apply(lines)
	if len(got) != 2 {
		t.Fatalf("expected 2 lines, got %d: %v", len(got), got)
	}
}

func TestClamper_Apply_OnlyMin(t *testing.T) {
	c, _ := New(Options{Pattern: `val=(\d+)`, Min: fptr(50)})
	lines := []string{"val=10", "val=50", "val=999"}
	got := c.Apply(lines)
	if len(got) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(got))
	}
}

func TestClamper_Apply_EmptyInput(t *testing.T) {
	c, _ := New(Options{Pattern: `n=(\d+)`, Min: fptr(0), Max: fptr(10)})
	got := c.Apply([]string{})
	if len(got) != 0 {
		t.Fatalf("expected empty output, got %d lines", len(got))
	}
}

func TestClamper_Apply_NoMatch(t *testing.T) {
	c, _ := New(Options{Pattern: `n=(\d+)`, Min: fptr(0), Max: fptr(10)})
	lines := []string{"no numbers here", "also nothing"}
	got := c.Apply(lines)
	if len(got) != 0 {
		t.Fatalf("expected 0 lines, got %d", len(got))
	}
}

func TestClamper_FormatSummary(t *testing.T) {
	c, _ := New(Options{Pattern: `n=(\d+)`, Min: fptr(1), Max: fptr(99)})
	s := c.FormatSummary(3, 10)
	if s == "" {
		t.Fatal("expected non-empty summary")
	}
}
