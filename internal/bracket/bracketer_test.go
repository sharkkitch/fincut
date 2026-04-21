package bracket

import (
	"strings"
	"testing"
)

func TestNew_MissingPattern(t *testing.T) {
	_, err := New(Options{Open: "[", Close: "]"})
	if err == nil {
		t.Fatal("expected error for missing pattern")
	}
}

func TestNew_EmptyBrackets(t *testing.T) {
	_, err := New(Options{WrapAll: true})
	if err == nil {
		t.Fatal("expected error when both Open and Close are empty")
	}
}

func TestNew_InvalidPattern(t *testing.T) {
	_, err := New(Options{Pattern: "[invalid", Open: "[", Close: "]"})
	if err == nil {
		t.Fatal("expected error for invalid regex")
	}
}

func TestNew_ValidOptions(t *testing.T) {
	b, err := New(Options{Pattern: `\d+`, Open: "<", Close: ">"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if b == nil {
		t.Fatal("expected non-nil Bracketer")
	}
}

func TestBracketer_Apply_MatchingLines(t *testing.T) {
	b, _ := New(Options{Pattern: `ERROR`, Open: ">>>", Close: "<<<"})
	input := []string{"INFO ok", "ERROR bad", "WARN maybe"}
	out := b.Apply(input)
	if out[0] != "INFO ok" {
		t.Errorf("non-matching line modified: %q", out[0])
	}
	if out[1] != ">>>ERROR bad<<<" {
		t.Errorf("expected wrapped line, got %q", out[1])
	}
	if out[2] != "WARN maybe" {
		t.Errorf("non-matching line modified: %q", out[2])
	}
}

func TestBracketer_Apply_WrapAll(t *testing.T) {
	b, _ := New(Options{WrapAll: true, Open: "[", Close: "]"})
	input := []string{"alpha", "beta", "gamma"}
	out := b.Apply(input)
	for i, l := range out {
		if !strings.HasPrefix(l, "[") || !strings.HasSuffix(l, "]") {
			t.Errorf("line %d not wrapped: %q", i, l)
		}
	}
}

func TestBracketer_Apply_EmptyInput(t *testing.T) {
	b, _ := New(Options{WrapAll: true, Open: "(", Close: ")"})
	out := b.Apply([]string{})
	if len(out) != 0 {
		t.Errorf("expected empty output, got %d lines", len(out))
	}
}

func TestBracketer_Apply_OriginalUnmodified(t *testing.T) {
	b, _ := New(Options{WrapAll: true, Open: "<", Close: ">"})
	input := []string{"hello", "world"}
	copy := append([]string{}, input...)
	b.Apply(input)
	for i := range input {
		if input[i] != copy[i] {
			t.Errorf("original slice modified at index %d", i)
		}
	}
}

func TestFormatSummary(t *testing.T) {
	s := FormatSummary(10, 4)
	if !strings.Contains(s, "4/10") {
		t.Errorf("unexpected summary format: %q", s)
	}
}
