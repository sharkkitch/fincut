package highlight

import (
	"strings"
	"testing"
)

func TestNew_NoPatterns(t *testing.T) {
	_, err := New(Options{Patterns: nil})
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

func TestNew_ValidPatterns(t *testing.T) {
	h, err := New(Options{Patterns: []string{"error", "warn"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if h == nil {
		t.Fatal("expected non-nil highlighter")
	}
}

func TestHighlighter_Apply_ContainsANSI(t *testing.T) {
	h, _ := New(Options{Patterns: []string{"error"}})
	result := h.Apply("this is an error message")
	if !strings.Contains(result, "\033[") {
		t.Errorf("expected ANSI codes in output, got: %q", result)
	}
}

func TestHighlighter_Apply_NoMatch(t *testing.T) {
	h, _ := New(Options{Patterns: []string{"error"}})
	line := "everything is fine"
	result := h.Apply(line)
	if result != line {
		t.Errorf("expected unchanged line, got: %q", result)
	}
}

func TestHighlighter_Apply_Bold(t *testing.T) {
	h, _ := New(Options{Patterns: []string{"warn"}, Bold: true})
	result := h.Apply("warn: disk usage high")
	if !strings.Contains(result, ansiBold) {
		t.Errorf("expected bold ANSI code in output, got: %q", result)
	}
}

func TestHighlighter_ApplyAll(t *testing.T) {
	h, _ := New(Options{Patterns: []string{"err"}})
	lines := []string{"ok line", "err: something failed", "another ok"}
	out := h.ApplyAll(lines)
	if len(out) != len(lines) {
		t.Fatalf("expected %d lines, got %d", len(lines), len(out))
	}
	if strings.Contains(out[0], "\033[") {
		t.Error("expected no ANSI in non-matching line")
	}
	if !strings.Contains(out[1], "\033[") {
		t.Error("expected ANSI in matching line")
	}
}

func TestStripANSI(t *testing.T) {
	h, _ := New(Options{Patterns: []string{"error"}})
	highlighted := h.Apply("this is an error message")
	stripped := StripANSI(highlighted)
	if strings.Contains(stripped, "\033[") {
		t.Errorf("expected no ANSI codes after strip, got: %q", stripped)
	}
	if !strings.Contains(stripped, "error") {
		t.Errorf("expected original text preserved, got: %q", stripped)
	}
}
