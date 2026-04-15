package tokenize

import (
	"testing"
)

func TestNew_NeitherSet(t *testing.T) {
	_, err := New(Options{})
	if err == nil {
		t.Fatal("expected error when neither Delimiter nor Pattern is set")
	}
}

func TestNew_MutuallyExclusive(t *testing.T) {
	_, err := New(Options{Delimiter: ",", Pattern: `(\w+)`})
	if err == nil {
		t.Fatal("expected error when both Delimiter and Pattern are set")
	}
}

func TestNew_InvalidPattern(t *testing.T) {
	_, err := New(Options{Pattern: `(unclosed`})
	if err == nil {
		t.Fatal("expected error for invalid regex pattern")
	}
}

func TestNew_PatternMissingCaptureGroup(t *testing.T) {
	_, err := New(Options{Pattern: `\w+`})
	if err == nil {
		t.Fatal("expected error when pattern has no capture groups")
	}
}

func TestNew_NegativeMinTokens(t *testing.T) {
	_, err := New(Options{Delimiter: ",", MinTokens: -1})
	if err == nil {
		t.Fatal("expected error for negative MinTokens")
	}
}

func TestNew_ValidDelimiter(t *testing.T) {
	_, err := New(Options{Delimiter: ","})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNew_ValidPattern(t *testing.T) {
	_, err := New(Options{Pattern: `(\w+)`})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestTokenizer_Apply_Delimiter(t *testing.T) {
	tok, _ := New(Options{Delimiter: ",", Join: "-"})
	got := tok.Apply([]string{"a,b,c", "x,y"})
	if len(got) != 2 || got[0] != "a-b-c" || got[1] != "x-y" {
		t.Fatalf("unexpected output: %v", got)
	}
}

func TestTokenizer_Apply_Pattern(t *testing.T) {
	tok, _ := New(Options{Pattern: `(\w+)\s+(\w+)`, Join: "|"})
	got := tok.Apply([]string{"hello world", "foo bar"})
	if len(got) != 2 || got[0] != "hello|world" || got[1] != "foo|bar" {
		t.Fatalf("unexpected output: %v", got)
	}
}

func TestTokenizer_Apply_MinTokens_DropsLine(t *testing.T) {
	tok, _ := New(Options{Delimiter: ",", MinTokens: 3})
	got := tok.Apply([]string{"a,b,c", "x,y"})
	if len(got) != 1 || got[0] != "a b c" {
		t.Fatalf("expected only the 3-token line, got: %v", got)
	}
}

func TestTokenizer_Apply_NoMatch_PatternDropsLine(t *testing.T) {
	tok, _ := New(Options{Pattern: `(\d+)-(\d+)`})
	got := tok.Apply([]string{"no-numbers-here", "12-34"})
	if len(got) != 1 || got[0] != "12 34" {
		t.Fatalf("unexpected output: %v", got)
	}
}

func TestTokenizer_Apply_EmptyInput(t *testing.T) {
	tok, _ := New(Options{Delimiter: " "})
	got := tok.Apply([]string{})
	if len(got) != 0 {
		t.Fatalf("expected empty output, got: %v", got)
	}
}
