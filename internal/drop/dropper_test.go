package drop

import (
	"strings"
	"testing"
)

func TestNew_NoPatterns(t *testing.T) {
	_, err := New(Options{})
	if err == nil {
		t.Fatal("expected error for empty patterns, got nil")
	}
}

func TestNew_InvalidPattern(t *testing.T) {
	_, err := New(Options{Patterns: []string{`[invalid`}})
	if err == nil {
		t.Fatal("expected error for invalid regex, got nil")
	}
}

func TestNew_ValidOptions(t *testing.T) {
	d, err := New(Options{Patterns: []string{`^DEBUG`}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if d == nil {
		t.Fatal("expected non-nil Dropper")
	}
}

func TestDropper_Apply_RemovesMatchingLines(t *testing.T) {
	d, _ := New(Options{Patterns: []string{`^DEBUG`}})
	input := []string{"INFO hello", "DEBUG verbose", "ERROR oops", "DEBUG again"}
	got := d.Apply(input)
	for _, line := range got {
		if strings.HasPrefix(line, "DEBUG") {
			t.Errorf("expected DEBUG line to be dropped, got %q", line)
		}
	}
	if len(got) != 2 {
		t.Errorf("expected 2 lines, got %d", len(got))
	}
}

func TestDropper_Apply_AllKeptWhenNoMatch(t *testing.T) {
	d, _ := New(Options{Patterns: []string{`^TRACE`}})
	input := []string{"INFO a", "WARN b", "ERROR c"}
	got := d.Apply(input)
	if len(got) != len(input) {
		t.Errorf("expected %d lines, got %d", len(input), len(got))
	}
}

func TestDropper_Apply_EmptyInput(t *testing.T) {
	d, _ := New(Options{Patterns: []string{`.*`}})
	got := d.Apply([]string{})
	if len(got) != 0 {
		t.Errorf("expected empty result, got %d lines", len(got))
	}
}

func TestDropper_Apply_Invert(t *testing.T) {
	d, _ := New(Options{Patterns: []string{`^ERROR`}, Invert: true})
	input := []string{"INFO ok", "ERROR bad", "WARN maybe", "ERROR worse"}
	got := d.Apply(input)
	if len(got) != 2 {
		t.Errorf("expected 2 lines (only ERRORs), got %d", len(got))
	}
	for _, line := range got {
		if !strings.HasPrefix(line, "ERROR") {
			t.Errorf("expected only ERROR lines in inverted mode, got %q", line)
		}
	}
}

func TestDropper_Apply_MultiplePatterns(t *testing.T) {
	d, _ := New(Options{Patterns: []string{`^DEBUG`, `^TRACE`}})
	input := []string{"INFO a", "DEBUG b", "TRACE c", "ERROR d"}
	got := d.Apply(input)
	if len(got) != 2 {
		t.Errorf("expected 2 lines, got %d", len(got))
	}
}

func TestFormatSummary(t *testing.T) {
	s := FormatSummary(10, 3)
	if !strings.Contains(s, "7/10") {
		t.Errorf("expected summary to contain '7/10', got %q", s)
	}
	if !strings.Contains(s, "3 dropped") {
		t.Errorf("expected summary to mention '3 dropped', got %q", s)
	}
}
