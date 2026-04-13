package indent

import (
	"strings"
	"testing"
)

func TestNew_NoOpConfig(t *testing.T) {
	_, err := New(Options{Depth: 0, StripExisting: false})
	if err == nil {
		t.Fatal("expected error for no-op config, got nil")
	}
}

func TestNew_ValidOptions(t *testing.T) {
	_, err := New(Options{Depth: 2})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNew_DefaultUnit(t *testing.T) {
	in, err := New(Options{Depth: 1})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if in.opts.Unit != "  " {
		t.Errorf("expected default unit '  ', got %q", in.opts.Unit)
	}
}

func TestIndenter_Apply_Indent(t *testing.T) {
	in, _ := New(Options{Depth: 1, Unit: "  "})
	lines := []string{"hello", "world"}
	got := in.Apply(lines)
	for _, l := range got {
		if !strings.HasPrefix(l, "  ") {
			t.Errorf("expected indented line, got %q", l)
		}
	}
}

func TestIndenter_Apply_Dedent(t *testing.T) {
	in, _ := New(Options{Depth: -1, Unit: "  "})
	lines := []string{"  hello", "  world", "nospace"}
	got := in.Apply(lines)
	if got[0] != "hello" {
		t.Errorf("expected 'hello', got %q", got[0])
	}
	if got[2] != "nospace" {
		t.Errorf("expected 'nospace' unchanged, got %q", got[2])
	}
}

func TestIndenter_Apply_StripExisting(t *testing.T) {
	in, _ := New(Options{Depth: 1, Unit: ">", StripExisting: true})
	lines := []string{"   leading spaces"}
	got := in.Apply(lines)
	if got[0] != ">leading spaces" {
		t.Errorf("expected '>leading spaces', got %q", got[0])
	}
}

func TestIndenter_Apply_StripOnly(t *testing.T) {
	in, _ := New(Options{Depth: 0, StripExisting: true, Unit: "  "})
	lines := []string{"\t  indented", "clean"}
	got := in.Apply(lines)
	if got[0] != "indented" {
		t.Errorf("expected 'indented', got %q", got[0])
	}
	if got[1] != "clean" {
		t.Errorf("expected 'clean', got %q", got[1])
	}
}

func TestFormatSummary_Indent(t *testing.T) {
	summary := FormatSummary(10, Options{Depth: 2, Unit: "  "})
	if !strings.Contains(summary, "indented") {
		t.Errorf("expected 'indented' in summary, got %q", summary)
	}
}

func TestFormatSummary_Dedent(t *testing.T) {
	summary := FormatSummary(5, Options{Depth: -1, Unit: "\t"})
	if !strings.Contains(summary, "dedented") {
		t.Errorf("expected 'dedented' in summary, got %q", summary)
	}
}
