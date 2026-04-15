package fold

import (
	"testing"
)

func TestNew_EmptyPattern(t *testing.T) {
	_, err := New(Options{})
	if err == nil {
		t.Fatal("expected error for empty pattern")
	}
}

func TestNew_InvalidPattern(t *testing.T) {
	_, err := New(Options{ContinuationPattern: "[invalid"})
	if err == nil {
		t.Fatal("expected error for invalid regex")
	}
}

func TestNew_ValidOptions(t *testing.T) {
	f, err := New(Options{ContinuationPattern: `^\s+`})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f == nil {
		t.Fatal("expected non-nil Folder")
	}
}

func TestNew_DefaultSeparator(t *testing.T) {
	f, err := New(Options{ContinuationPattern: `^\s+`})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.separator != " " {
		t.Errorf("expected default separator ' ', got %q", f.separator)
	}
}

func TestFolder_Apply_EmptyInput(t *testing.T) {
	f, _ := New(Options{ContinuationPattern: `^\s+`})
	out := f.Apply([]string{})
	if len(out) != 0 {
		t.Errorf("expected empty output, got %v", out)
	}
}

func TestFolder_Apply_NoContinuation(t *testing.T) {
	f, _ := New(Options{ContinuationPattern: `^\s+`})
	in := []string{"line1", "line2", "line3"}
	out := f.Apply(in)
	if len(out) != 3 {
		t.Errorf("expected 3 lines, got %d", len(out))
	}
}

func TestFolder_Apply_MergesContinuation(t *testing.T) {
	f, _ := New(Options{ContinuationPattern: `^\s+`, Separator: " "})
	in := []string{"first", "  second", "  third", "fourth"}
	out := f.Apply(in)
	if len(out) != 2 {
		t.Fatalf("expected 2 lines, got %d: %v", len(out), out)
	}
	if out[0] != "first   second   third" {
		t.Errorf("unexpected merged line: %q", out[0])
	}
	if out[1] != "fourth" {
		t.Errorf("unexpected second line: %q", out[1])
	}
}

func TestFolder_Apply_CustomSeparator(t *testing.T) {
	f, _ := New(Options{ContinuationPattern: `\\$`, Separator: ""})
	in := []string{`hello\`, `world`}
	out := f.Apply(in)
	if len(out) != 1 {
		t.Fatalf("expected 1 line, got %d", len(out))
	}
	if out[0] != `hello\world` {
		t.Errorf("unexpected merged line: %q", out[0])
	}
}

func TestFormatSummary(t *testing.T) {
	s := FormatSummary(10, 7)
	if s == "" {
		t.Error("expected non-empty summary")
	}
}
