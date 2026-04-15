package context

import (
	"strings"
	"testing"
)

func TestNew_EmptyPattern(t *testing.T) {
	_, err := New(Options{Pattern: "", Before: 1, After: 1})
	if err == nil {
		t.Fatal("expected error for empty pattern")
	}
}

func TestNew_NegativeBefore(t *testing.T) {
	_, err := New(Options{Pattern: "foo", Before: -1, After: 0})
	if err == nil {
		t.Fatal("expected error for negative before")
	}
}

func TestNew_NegativeAfter(t *testing.T) {
	_, err := New(Options{Pattern: "foo", Before: 0, After: -2})
	if err == nil {
		t.Fatal("expected error for negative after")
	}
}

func TestNew_InvalidPattern(t *testing.T) {
	_, err := New(Options{Pattern: "[", Before: 0, After: 0})
	if err == nil {
		t.Fatal("expected error for invalid regex")
	}
}

func TestNew_ValidOptions(t *testing.T) {
	c, err := New(Options{Pattern: "error", Before: 2, After: 2})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c == nil {
		t.Fatal("expected non-nil Contexter")
	}
}

func TestContexter_Apply_NoMatch(t *testing.T) {
	c, _ := New(Options{Pattern: "ERROR", Before: 1, After: 1})
	lines := []string{"info: start", "info: running", "info: done"}
	matches := c.Apply(lines)
	if len(matches) != 0 {
		t.Fatalf("expected 0 matches, got %d", len(matches))
	}
}

func TestContexter_Apply_BasicMatch(t *testing.T) {
	c, _ := New(Options{Pattern: "ERROR", Before: 1, After: 1})
	lines := []string{"line1", "ERROR here", "line3"}
	matches := c.Apply(lines)
	if len(matches) != 1 {
		t.Fatalf("expected 1 match, got %d", len(matches))
	}
	if matches[0].LineNo != 2 {
		t.Errorf("expected LineNo=2, got %d", matches[0].LineNo)
	}
	if len(matches[0].Lines) != 3 {
		t.Errorf("expected 3 context lines, got %d", len(matches[0].Lines))
	}
}

func TestContexter_Apply_ClampsBounds(t *testing.T) {
	c, _ := New(Options{Pattern: "hit", Before: 5, After: 5})
	lines := []string{"hit", "b", "c"}
	matches := c.Apply(lines)
	if len(matches) != 1 {
		t.Fatalf("expected 1 match")
	}
	if len(matches[0].Lines) != 3 {
		t.Errorf("expected 3 lines (clamped), got %d", len(matches[0].Lines))
	}
}

func TestContexter_Apply_Label(t *testing.T) {
	c, _ := New(Options{Pattern: "target", Before: 0, After: 0, Label: true})
	lines := []string{"alpha", "target line", "beta"}
	matches := c.Apply(lines)
	if len(matches) != 1 {
		t.Fatalf("expected 1 match")
	}
	if !strings.HasPrefix(matches[0].Lines[0], "2: ") {
		t.Errorf("expected label prefix '2: ', got %q", matches[0].Lines[0])
	}
}

func TestContexter_Apply_MultipleMatches(t *testing.T) {
	c, _ := New(Options{Pattern: "WARN", Before: 0, After: 1})
	lines := []string{"ok", "WARN one", "detail", "ok", "WARN two", "end"}
	matches := c.Apply(lines)
	if len(matches) != 2 {
		t.Fatalf("expected 2 matches, got %d", len(matches))
	}
	if matches[1].LineNo != 5 {
		t.Errorf("expected second match at line 5, got %d", matches[1].LineNo)
	}
}
