package join

import (
	"testing"
)

func TestNew_NegativeGroupSize(t *testing.T) {
	_, err := New(Options{GroupSize: -1})
	if err == nil {
		t.Fatal("expected error for negative GroupSize")
	}
}

func TestNew_ValidOptions(t *testing.T) {
	j, err := New(Options{GroupSize: 2, Separator: ","})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if j == nil {
		t.Fatal("expected non-nil Joiner")
	}
}

func TestNew_DefaultSeparator(t *testing.T) {
	j, err := New(Options{GroupSize: 2})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if j.opts.Separator != " " {
		t.Errorf("expected default separator ' ', got %q", j.opts.Separator)
	}
}

func TestJoiner_Apply_EmptyInput(t *testing.T) {
	j, _ := New(Options{GroupSize: 2})
	out := j.Apply([]string{})
	if len(out) != 0 {
		t.Errorf("expected 0 lines, got %d", len(out))
	}
}

func TestJoiner_Apply_GroupSize2(t *testing.T) {
	j, _ := New(Options{GroupSize: 2, Separator: "|" })
	input := []string{"a", "b", "c", "d"}
	out := j.Apply(input)
	if len(out) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(out))
	}
	if out[0] != "a|b" {
		t.Errorf("expected 'a|b', got %q", out[0])
	}
	if out[1] != "c|d" {
		t.Errorf("expected 'c|d', got %q", out[1])
	}
}

func TestJoiner_Apply_OddLines(t *testing.T) {
	j, _ := New(Options{GroupSize: 2, Separator: "-"})
	input := []string{"x", "y", "z"}
	out := j.Apply(input)
	if len(out) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(out))
	}
	if out[1] != "z" {
		t.Errorf("expected last group to be 'z', got %q", out[1])
	}
}

func TestJoiner_Apply_GroupSizeZero_AllJoined(t *testing.T) {
	j, _ := New(Options{GroupSize: 0, Separator: ", "})
	input := []string{"one", "two", "three"}
	out := j.Apply(input)
	if len(out) != 1 {
		t.Fatalf("expected 1 line, got %d", len(out))
	}
	if out[0] != "one, two, three" {
		t.Errorf("unexpected output: %q", out[0])
	}
}

func TestFormatSummary(t *testing.T) {
	s := FormatSummary(10, 5)
	if s == "" {
		t.Error("expected non-empty summary")
	}
}
