package prefix

import (
	"strings"
	"testing"
)

func TestNew_NeitherSet(t *testing.T) {
	_, err := New(Options{})
	if err == nil {
		t.Fatal("expected error when neither Text nor LineNumbers is set")
	}
}

func TestNew_MutuallyExclusive(t *testing.T) {
	_, err := New(Options{Text: "LOG", LineNumbers: true})
	if err == nil {
		t.Fatal("expected error when both Text and LineNumbers are set")
	}
}

func TestNew_NegativeWidth(t *testing.T) {
	_, err := New(Options{LineNumbers: true, Width: -1})
	if err == nil {
		t.Fatal("expected error for negative Width")
	}
}

func TestNew_ValidText(t *testing.T) {
	p, err := New(Options{Text: ">>>"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p == nil {
		t.Fatal("expected non-nil Prefixer")
	}
}

func TestPrefixer_Apply_TextPrefix(t *testing.T) {
	p, _ := New(Options{Text: "LOG", Separator: " | "})
	input := []string{"hello", "world"}
	got := p.Apply(input)
	if len(got) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(got))
	}
	if got[0] != "LOG | hello" {
		t.Errorf("line 0: got %q", got[0])
	}
	if got[1] != "LOG | world" {
		t.Errorf("line 1: got %q", got[1])
	}
}

func TestPrefixer_Apply_LineNumbers(t *testing.T) {
	p, _ := New(Options{LineNumbers: true})
	input := []string{"alpha", "beta", "gamma"}
	got := p.Apply(input)
	for i, line := range got {
		if !strings.HasPrefix(line, strings.Repeat("", 0)) {
			_ = i
		}
	}
	if !strings.HasPrefix(got[0], "1: ") {
		t.Errorf("expected '1: alpha', got %q", got[0])
	}
	if !strings.HasPrefix(got[2], "3: ") {
		t.Errorf("expected '3: gamma', got %q", got[2])
	}
}

func TestPrefixer_Apply_LineNumbers_ZeroPad(t *testing.T) {
	p, _ := New(Options{LineNumbers: true, Width: 3})
	got := p.Apply([]string{"x"})
	if got[0] != "001: x" {
		t.Errorf("expected '001: x', got %q", got[0])
	}
}

func TestPrefixer_Apply_DefaultSeparator(t *testing.T) {
	p, _ := New(Options{Text: "T"})
	got := p.Apply([]string{"line"})
	if got[0] != "T: line" {
		t.Errorf("expected 'T: line', got %q", got[0])
	}
}

func TestFormatSummary_Text(t *testing.T) {
	opts := Options{Text: "ERR", Separator: ": "}
	s := FormatSummary(5, opts)
	if !strings.Contains(s, "5") || !strings.Contains(s, "ERR") {
		t.Errorf("unexpected summary: %q", s)
	}
}

func TestFormatSummary_LineNumbers(t *testing.T) {
	opts := Options{LineNumbers: true, Separator: ": "}
	s := FormatSummary(10, opts)
	if !strings.Contains(s, "10") || !strings.Contains(s, "line numbers") {
		t.Errorf("unexpected summary: %q", s)
	}
}
