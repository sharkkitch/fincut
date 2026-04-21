package center

import (
	"strings"
	"testing"
)

func TestNew_ZeroWidth(t *testing.T) {
	_, err := New(Options{Width: 0})
	if err == nil {
		t.Fatal("expected error for zero width")
	}
}

func TestNew_NegativeWidth(t *testing.T) {
	_, err := New(Options{Width: -5})
	if err == nil {
		t.Fatal("expected error for negative width")
	}
}

func TestNew_ValidOptions(t *testing.T) {
	c, err := New(Options{Width: 20})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c == nil {
		t.Fatal("expected non-nil centerer")
	}
}

func TestNew_DefaultFill(t *testing.T) {
	c, err := New(Options{Width: 10})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.fill != " " {
		t.Errorf("expected default fill ' ', got %q", c.fill)
	}
}

func TestCenterer_Apply_ShortLine(t *testing.T) {
	c, _ := New(Options{Width: 10})
	out := c.Apply([]string{"hi"})
	if len(out) != 1 {
		t.Fatalf("expected 1 line, got %d", len(out))
	}
	if len(out[0]) != 10 {
		t.Errorf("expected width 10, got %d: %q", len(out[0]), out[0])
	}
	if !strings.Contains(out[0], "hi") {
		t.Errorf("output does not contain original text: %q", out[0])
	}
}

func TestCenterer_Apply_ExactWidth(t *testing.T) {
	c, _ := New(Options{Width: 4})
	out := c.Apply([]string{"abcd"})
	if out[0] != "abcd" {
		t.Errorf("expected unchanged line, got %q", out[0])
	}
}

func TestCenterer_Apply_LongerThanWidth(t *testing.T) {
	c, _ := New(Options{Width: 3})
	out := c.Apply([]string{"toolong"})
	if out[0] != "toolong" {
		t.Errorf("expected unchanged line, got %q", out[0])
	}
}

func TestCenterer_Apply_EmptyInput(t *testing.T) {
	c, _ := New(Options{Width: 10})
	out := c.Apply([]string{})
	if len(out) != 0 {
		t.Errorf("expected empty output, got %d lines", len(out))
	}
}

func TestCenterer_Apply_CustomFill(t *testing.T) {
	c, _ := New(Options{Width: 9, Fill: "-"})
	out := c.Apply([]string{"hi"})
	if out[0] != "---hi----" && out[0] != "----hi---" {
		t.Errorf("unexpected centered output: %q", out[0])
	}
}

func TestFormatSummary(t *testing.T) {
	c, _ := New(Options{Width: 80})
	s := FormatSummary(c)
	if !strings.Contains(s, "80") {
		t.Errorf("summary missing width: %q", s)
	}
}
