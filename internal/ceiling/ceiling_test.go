package ceiling

import (
	"strings"
	"testing"
)

func TestNew_EmptyPattern(t *testing.T) {
	_, err := New(Options{Pattern: "", Max: 3})
	if err == nil {
		t.Fatal("expected error for empty pattern")
	}
}

func TestNew_ZeroMax(t *testing.T) {
	_, err := New(Options{Pattern: "err", Max: 0})
	if err == nil {
		t.Fatal("expected error for zero max")
	}
}

func TestNew_NegativeMax(t *testing.T) {
	_, err := New(Options{Pattern: "err", Max: -1})
	if err == nil {
		t.Fatal("expected error for negative max")
	}
}

func TestNew_InvalidPattern(t *testing.T) {
	_, err := New(Options{Pattern: "[", Max: 1})
	if err == nil {
		t.Fatal("expected error for invalid regex")
	}
}

func TestNew_ValidOptions(t *testing.T) {
	c, err := New(Options{Pattern: "ERROR", Max: 5})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c == nil {
		t.Fatal("expected non-nil Ceiling")
	}
}

func TestCeiling_Apply_NonMatchingPassThrough(t *testing.T) {
	c, _ := New(Options{Pattern: "ERROR", Max: 2})
	lines := []string{"info a", "info b", "info c"}
	out := c.Apply(lines)
	if len(out) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(out))
	}
}

func TestCeiling_Apply_DrAfterMax(t *testing.T) {
	c, _ := New(Options{Pattern: "ERROR", Max: 2})
	lines := []string{"ERROR 1", "ERROR 2", "ERROR 3", "ERROR 4"}
	out := c.Apply(lines)
	if len(out) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(out))
	}
}

func TestCeiling_Apply_MixedLines(t *testing.T) {
	c, _ := New(Options{Pattern: "ERROR", Max: 2})
	lines := []string{"INFO start", "ERROR 1", "INFO mid", "ERROR 2", "ERROR 3", "INFO end"}
	out := c.Apply(lines)
	// expect: INFO start, ERROR 1, INFO mid, ERROR 2, INFO end (ERROR 3 dropped)
	if len(out) != 5 {
		t.Fatalf("expected 5 lines, got %d: %v", len(out), out)
	}
	for _, l := range out {
		if l == "ERROR 3" {
			t.Fatal("ERROR 3 should have been dropped")
		}
	}
}

func TestCeiling_FormatSummary(t *testing.T) {
	c, _ := New(Options{Pattern: "WARN", Max: 10})
	c.Apply([]string{"WARN a", "WARN b"})
	s := c.FormatSummary()
	if !strings.Contains(s, "2/10") {
		t.Fatalf("unexpected summary: %s", s)
	}
}
