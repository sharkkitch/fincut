package between

import (
	"testing"
)

func TestNew_EmptyStartPattern(t *testing.T) {
	_, err := New(Options{StartPattern: "", EndPattern: "END"})
	if err == nil {
		t.Fatal("expected error for empty start pattern")
	}
}

func TestNew_EmptyEndPattern(t *testing.T) {
	_, err := New(Options{StartPattern: "START", EndPattern: ""})
	if err == nil {
		t.Fatal("expected error for empty end pattern")
	}
}

func TestNew_InvalidStartPattern(t *testing.T) {
	_, err := New(Options{StartPattern: "[", EndPattern: "END"})
	if err == nil {
		t.Fatal("expected error for invalid start pattern")
	}
}

func TestNew_ValidOptions(t *testing.T) {
	_, err := New(Options{StartPattern: "START", EndPattern: "END"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestBetweener_Apply_BasicExtract(t *testing.T) {
	b, _ := New(Options{StartPattern: "START", EndPattern: "END"})
	lines := []string{"before", "START", "inside1", "inside2", "END", "after"}
	out := b.Apply(lines)
	if len(out) != 2 || out[0] != "inside1" || out[1] != "inside2" {
		t.Fatalf("unexpected output: %v", out)
	}
}

func TestBetweener_Apply_Inclusive(t *testing.T) {
	b, _ := New(Options{StartPattern: "START", EndPattern: "END", Inclusive: true})
	lines := []string{"START", "middle", "END"}
	out := b.Apply(lines)
	if len(out) != 3 {
		t.Fatalf("expected 3 lines, got %d: %v", len(out), out)
	}
}

func TestBetweener_Apply_MultipleRegions(t *testing.T) {
	b, _ := New(Options{StartPattern: "START", EndPattern: "END"})
	lines := []string{"START", "a", "END", "gap", "START", "b", "END"}
	out := b.Apply(lines)
	if len(out) != 2 || out[0] != "a" || out[1] != "b" {
		t.Fatalf("unexpected output: %v", out)
	}
}

func TestBetweener_Apply_NoMatch(t *testing.T) {
	b, _ := New(Options{StartPattern: "START", EndPattern: "END"})
	lines := []string{"foo", "bar", "baz"}
	out := b.Apply(lines)
	if len(out) != 0 {
		t.Fatalf("expected empty output, got %v", out)
	}
}

func TestBetweener_Apply_EmptyInput(t *testing.T) {
	b, _ := New(Options{StartPattern: "START", EndPattern: "END"})
	out := b.Apply([]string{})
	if len(out) != 0 {
		t.Fatalf("expected empty output, got %v", out)
	}
}

func TestFormatSummary(t *testing.T) {
	s := FormatSummary(100, 42)
	if s == "" {
		t.Fatal("expected non-empty summary")
	}
}
