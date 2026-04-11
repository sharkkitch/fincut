package grep

import (
	"testing"
)

func TestNew_NoPatterns(t *testing.T) {
	_, err := New(Options{})
	if err == nil {
		t.Fatal("expected error for empty patterns")
	}
}

func TestNew_InvalidPattern(t *testing.T) {
	_, err := New(Options{Patterns: []string{"[invalid"}})
	if err == nil {
		t.Fatal("expected error for invalid regex")
	}
}

func TestNew_NegativeContext(t *testing.T) {
	_, err := New(Options{Patterns: []string{"foo"}, ContextBefore: -1})
	if err == nil {
		t.Fatal("expected error for negative context")
	}
}

func TestNew_ValidOptions(t *testing.T) {
	g, err := New(Options{Patterns: []string{`\d+`}, ContextBefore: 1, ContextAfter: 2})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if g == nil {
		t.Fatal("expected non-nil grepper")
	}
}

func TestGrepper_Apply_BasicMatch(t *testing.T) {
	g, _ := New(Options{Patterns: []string{"ERROR"}})
	lines := []string{"INFO ok", "ERROR boom", "INFO fine"}
	matches := g.Apply(lines)
	if len(matches) != 1 {
		t.Fatalf("expected 1 match, got %d", len(matches))
	}
	if matches[0].LineNumber != 2 {
		t.Errorf("expected line 2, got %d", matches[0].LineNumber)
	}
}

func TestGrepper_Apply_Invert(t *testing.T) {
	g, _ := New(Options{Patterns: []string{"ERROR"}, Invert: true})
	lines := []string{"INFO ok", "ERROR boom", "INFO fine"}
	matches := g.Apply(lines)
	if len(matches) != 2 {
		t.Fatalf("expected 2 matches, got %d", len(matches))
	}
}

func TestGrepper_Apply_ContextBefore(t *testing.T) {
	g, _ := New(Options{Patterns: []string{"ERROR"}, ContextBefore: 1})
	lines := []string{"INFO a", "INFO b", "ERROR c"}
	matches := g.Apply(lines)
	if len(matches[0].Before) != 1 || matches[0].Before[0] != "INFO b" {
		t.Errorf("unexpected before context: %v", matches[0].Before)
	}
}

func TestGrepper_Apply_ContextAfter(t *testing.T) {
	g, _ := New(Options{Patterns: []string{"ERROR"}, ContextAfter: 2})
	lines := []string{"ERROR x", "INFO y", "INFO z", "INFO w"}
	matches := g.Apply(lines)
	if len(matches[0].After) != 2 {
		t.Errorf("expected 2 after lines, got %d", len(matches[0].After))
	}
}

func TestGrepper_Apply_NoMatch(t *testing.T) {
	g, _ := New(Options{Patterns: []string{"FATAL"}})
	lines := []string{"INFO ok", "DEBUG trace"}
	matches := g.Apply(lines)
	if len(matches) != 0 {
		t.Errorf("expected no matches, got %d", len(matches))
	}
}

func TestGrepper_Apply_MultiPattern(t *testing.T) {
	g, _ := New(Options{Patterns: []string{"ERROR", "WARN"}})
	lines := []string{"INFO ok", "WARN slow", "ERROR fail", "DEBUG trace"}
	matches := g.Apply(lines)
	if len(matches) != 2 {
		t.Fatalf("expected 2 matches, got %d", len(matches))
	}
}
