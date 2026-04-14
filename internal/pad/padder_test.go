package pad

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
	p, err := New(Options{Width: 20})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p == nil {
		t.Fatal("expected non-nil Padder")
	}
}

func TestNew_DefaultFill(t *testing.T) {
	p, err := New(Options{Width: 10})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.opts.Fill != ' ' {
		t.Errorf("expected default fill ' ', got %q", p.opts.Fill)
	}
}

func TestPadder_Apply_PadsRight(t *testing.T) {
	p, _ := New(Options{Width: 10, Fill: ' '})
	out := p.Apply([]string{"hello"})
	if len(out) != 1 {
		t.Fatalf("expected 1 line, got %d", len(out))
	}
	if out[0] != "hello     " {
		t.Errorf("unexpected output: %q", out[0])
	}
}

func TestPadder_Apply_PadsLeft(t *testing.T) {
	p, _ := New(Options{Width: 10, Fill: ' ', Left: true})
	out := p.Apply([]string{"hello"})
	if out[0] != "     hello" {
		t.Errorf("unexpected output: %q", out[0])
	}
}

func TestPadder_Apply_NoOpWhenAtWidth(t *testing.T) {
	p, _ := New(Options{Width: 5, Fill: '-'})
	out := p.Apply([]string{"hello"})
	if out[0] != "hello" {
		t.Errorf("expected unchanged line, got %q", out[0])
	}
}

func TestPadder_Apply_NoOpWhenExceedsWidth(t *testing.T) {
	p, _ := New(Options{Width: 3, Fill: '-'})
	out := p.Apply([]string{"toolong"})
	if out[0] != "toolong" {
		t.Errorf("expected unchanged line, got %q", out[0])
	}
}

func TestPadder_Apply_CustomFill(t *testing.T) {
	p, _ := New(Options{Width: 8, Fill: '.'})
	out := p.Apply([]string{"go"})
	if out[0] != "go......" {
		t.Errorf("unexpected output: %q", out[0])
	}
}

func TestPadder_Apply_MultiLine(t *testing.T) {
	p, _ := New(Options{Width: 6, Fill: ' '})
	in := []string{"a", "bb", "ccc", "dddddd", "eeeeeee"}
	out := p.Apply(in)
	for i, line := range out {
		if len([]rune(line)) < 6 && i < 4 {
			t.Errorf("line %d too short: %q", i, line)
		}
		if !strings.HasPrefix(line, in[i]) {
			t.Errorf("line %d lost original content: %q", i, line)
		}
	}
}

func TestFormatSummary(t *testing.T) {
	s := FormatSummary(Options{Width: 80, Fill: '-', Left: true})
	if s == "" {
		t.Error("expected non-empty summary")
	}
}
