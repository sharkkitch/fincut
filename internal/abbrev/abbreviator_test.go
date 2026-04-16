package abbrev

import (
	"strings"
	"testing"
)

func TestNew_ZeroMaxTokenLen(t *testing.T) {
	_, err := New(Options{MaxTokenLen: 0, Delimiter: " "})
	if err == nil {
		t.Fatal("expected error for zero MaxTokenLen")
	}
}

func TestNew_NegativeMaxTokenLen(t *testing.T) {
	_, err := New(Options{MaxTokenLen: -1, Delimiter: " "})
	if err == nil {
		t.Fatal("expected error for negative MaxTokenLen")
	}
}

func TestNew_DefaultDelimiterAndEllipsis(t *testing.T) {
	a, err := New(Options{MaxTokenLen: 5})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if a.opts.Delimiter != " " {
		t.Errorf("expected default delimiter ' ', got %q", a.opts.Delimiter)
	}
	if a.opts.Ellipsis != "…" {
		t.Errorf("expected default ellipsis '…', got %q", a.opts.Ellipsis)
	}
}

func TestNew_ValidOptions(t *testing.T) {
	_, err := New(Options{MaxTokenLen: 3, Delimiter: ",", Ellipsis: "..."})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestAbbreviator_Apply_ShortTokensUnchanged(t *testing.T) {
	a, _ := New(Options{MaxTokenLen: 10, Delimiter: " "})
	input := []string{"hello world", "foo bar"}
	out := a.Apply(input)
	for i, line := range input {
		if out[i] != line {
			t.Errorf("line %d: expected %q, got %q", i, line, out[i])
		}
	}
}

func TestAbbreviator_Apply_LongTokenTruncated(t *testing.T) {
	a, _ := New(Options{MaxTokenLen: 4, Delimiter: " ", Ellipsis: "..."})
	out := a.Apply([]string{"hello world"})
	if !strings.HasPrefix(out[0], "hell...") {
		t.Errorf("expected truncated token, got %q", out[0])
	}
}

func TestAbbreviator_Apply_CustomDelimiter(t *testing.T) {
	a, _ := New(Options{MaxTokenLen: 3, Delimiter: ",", Ellipsis: "-"})
	out := a.Apply([]string{"abc,defgh,xy"})
	if out[0] != "abc,def-,xy" {
		t.Errorf("unexpected output: %q", out[0])
	}
}

func TestAbbreviator_Apply_EmptyInput(t *testing.T) {
	a, _ := New(Options{MaxTokenLen: 5})
	out := a.Apply([]string{})
	if len(out) != 0 {
		t.Errorf("expected empty output, got %d lines", len(out))
	}
}

func TestFormatSummary(t *testing.T) {
	in := []string{"hello world", "short"}
	out := []string{"hell… worl…", "short"}
	s := FormatSummary(in, out)
	if !strings.Contains(s, "1/2") {
		t.Errorf("expected '1/2' in summary, got %q", s)
	}
}
