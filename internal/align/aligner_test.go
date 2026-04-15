package align

import (
	"strings"
	"testing"
)

func TestNew_EmptyDelimiter(t *testing.T) {
	_, err := New(Options{Delimiter: ""})
	if err == nil {
		t.Fatal("expected error for empty delimiter")
	}
}

func TestNew_NegativePadding(t *testing.T) {
	_, err := New(Options{Delimiter: "|", Padding: -1})
	if err == nil {
		t.Fatal("expected error for negative padding")
	}
}

func TestNew_ValidOptions(t *testing.T) {
	a, err := New(Options{Delimiter: "|", Padding: 2})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if a == nil {
		t.Fatal("expected non-nil aligner")
	}
}

func TestNew_DefaultPadding(t *testing.T) {
	a, err := New(Options{Delimiter: ","})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if a.opts.Padding != 1 {
		t.Fatalf("expected default padding 1, got %d", a.opts.Padding)
	}
}

func TestAligner_Apply_EmptyInput(t *testing.T) {
	a, _ := New(Options{Delimiter: "|"})
	out := a.Apply(nil)
	if len(out) != 0 {
		t.Fatalf("expected empty output, got %d lines", len(out))
	}
}

func TestAligner_Apply_SingleLine(t *testing.T) {
	a, _ := New(Options{Delimiter: "|", Padding: 1})
	out := a.Apply([]string{"a|b|c"})
	if len(out) != 1 {
		t.Fatalf("expected 1 line, got %d", len(out))
	}
}

func TestAligner_Apply_ColumnsAligned(t *testing.T) {
	a, _ := New(Options{Delimiter: "|", Padding: 1})
	lines := []string{
		"foo|bar|baz",
		"x|longervalue|y",
	}
	out := a.Apply(lines)
	if len(out) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(out))
	}
	// Both lines should start with a field padded to len("foo") or len("x"), whichever is wider.
	if !strings.HasPrefix(out[1], "x  ") {
		t.Errorf("expected second line first field padded, got %q", out[1])
	}
}

func TestAligner_Apply_TrimFields(t *testing.T) {
	a, _ := New(Options{Delimiter: "|", Padding: 1, TrimFields: true})
	lines := []string{" foo | bar ", " x | y "}
	out := a.Apply(lines)
	for _, l := range out {
		if strings.Contains(l, "  ") && strings.HasPrefix(l, " ") {
			t.Errorf("line still has leading space after trim: %q", l)
		}
	}
}

func TestAligner_Apply_NoDelimiterInLine(t *testing.T) {
	a, _ := New(Options{Delimiter: "|", Padding: 1})
	lines := []string{"no-delimiter-here", "also-none"}
	out := a.Apply(lines)
	if out[0] != "no-delimiter-here" {
		t.Errorf("single-field line should be unchanged, got %q", out[0])
	}
}

func TestFormatSummary(t *testing.T) {
	summary := FormatSummary([]string{"a", "b"}, []string{"a", "b"})
	if !strings.Contains(summary, "2") {
		t.Errorf("expected summary to mention line count, got %q", summary)
	}
}
