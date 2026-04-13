package compact

import (
	"testing"
)

func TestNew_ValidOptions(t *testing.T) {
	c, err := New(Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c == nil {
		t.Fatal("expected non-nil Compacter")
	}
}

func TestCompacter_Apply_EmptyInput(t *testing.T) {
	c, _ := New(Options{})
	out, err := c.Apply([]string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 0 {
		t.Fatalf("expected empty output, got %d lines", len(out))
	}
}

func TestCompacter_Apply_RemovesBlankLines(t *testing.T) {
	c, _ := New(Options{})
	input := []string{"hello", "", "world", "", ""}
	out, err := c.Apply(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 2 {
		t.Fatalf("expected 2 lines, got %d: %v", len(out), out)
	}
}

func TestCompacter_Apply_CollapseBlank(t *testing.T) {
	c, _ := New(Options{CollapseBlank: true})
	input := []string{"a", "", "", "", "b"}
	out, err := c.Apply(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Expect: "a", "", "b"
	if len(out) != 3 {
		t.Fatalf("expected 3 lines, got %d: %v", len(out), out)
	}
	if out[1] != "" {
		t.Fatalf("expected blank middle line, got %q", out[1])
	}
}

func TestCompacter_Apply_TrimSpace(t *testing.T) {
	c, _ := New(Options{TrimSpace: true})
	input := []string{"  hello  ", "   ", "\t", "world"}
	out, err := c.Apply(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 2 {
		t.Fatalf("expected 2 lines, got %d: %v", len(out), out)
	}
	if out[0] != "hello" {
		t.Fatalf("expected trimmed line, got %q", out[0])
	}
}

func TestCompacter_Apply_AllBlank(t *testing.T) {
	c, _ := New(Options{})
	input := []string{"", "", ""}
	out, err := c.Apply(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 0 {
		t.Fatalf("expected empty output, got %d lines", len(out))
	}
}

func TestCompacter_Apply_NoBlankLines(t *testing.T) {
	c, _ := New(Options{CollapseBlank: true, TrimSpace: true})
	input := []string{"alpha", "beta", "gamma"}
	out, err := c.Apply(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(out))
	}
}
