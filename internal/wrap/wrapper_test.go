package wrap

import (
	"strings"
	"testing"
)

func TestNew_ZeroWidth(t *testing.T) {
	_, err := New(Options{Width: 0})
	if err == nil {
		t.Fatal("expected error for zero Width")
	}
}

func TestNew_NegativeWidth(t *testing.T) {
	_, err := New(Options{Width: -5})
	if err == nil {
		t.Fatal("expected error for negative Width")
	}
}

func TestNew_ValidOptions(t *testing.T) {
	w, err := New(Options{Width: 80})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if w == nil {
		t.Fatal("expected non-nil Wrapper")
	}
}

func TestWrapper_Apply_ShortLinesUnchanged(t *testing.T) {
	w, _ := New(Options{Width: 40})
	input := []string{"hello world", "foo bar"}
	out := w.Apply(input)
	if len(out) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(out))
	}
	for i, line := range out {
		if line != input[i] {
			t.Errorf("line %d changed: got %q", i, line)
		}
	}
}

func TestWrapper_Apply_HardBreak(t *testing.T) {
	w, _ := New(Options{Width: 10, HardBreak: true})
	input := []string{"abcdefghijklmnopqrst"}
	out := w.Apply(input)
	if len(out) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(out))
	}
	for _, line := range out {
		if len([]rune(line)) > 10 {
			t.Errorf("line exceeds width: %q", line)
		}
	}
}

func TestWrapper_Apply_SoftBreakAtSpace(t *testing.T) {
	w, _ := New(Options{Width: 10, HardBreak: false})
	// 12 chars but space at position 5
	input := []string{"hello world!"}
	out := w.Apply(input)
	if len(out) < 2 {
		t.Fatalf("expected at least 2 lines, got %d: %v", len(out), out)
	}
	if strings.Contains(out[0], " ") && len([]rune(out[0])) > 10 {
		t.Errorf("first segment too long: %q", out[0])
	}
}

func TestWrapper_Apply_Indent(t *testing.T) {
	w, _ := New(Options{Width: 10, Indent: "  ", HardBreak: true})
	input := []string{"abcdefghijklmnop"}
	out := w.Apply(input)
	if len(out) < 2 {
		t.Fatalf("expected wrapped output, got %d lines", len(out))
	}
	for _, line := range out[1:] {
		if !strings.HasPrefix(line, "  ") {
			t.Errorf("continuation line missing indent: %q", line)
		}
	}
}

func TestWrapper_Apply_EmptyInput(t *testing.T) {
	w, _ := New(Options{Width: 20})
	out := w.Apply([]string{})
	if len(out) != 0 {
		t.Fatalf("expected empty output, got %d lines", len(out))
	}
}

func TestFormatSummary(t *testing.T) {
	s := FormatSummary(3, 7)
	if !strings.Contains(s, "3") || !strings.Contains(s, "7") {
		t.Errorf("unexpected summary: %q", s)
	}
}
