package blank

import (
	"strings"
	"testing"
)

func TestNew_EveryTooSmall(t *testing.T) {
	_, err := New(Options{Every: 1})
	if err == nil {
		t.Fatal("expected error for Every < 2")
	}
}

func TestNew_NegativeOffset(t *testing.T) {
	_, err := New(Options{Every: 3, Offset: -1})
	if err == nil {
		t.Fatal("expected error for negative Offset")
	}
}

func TestNew_ValidOptions(t *testing.T) {
	b, err := New(Options{Every: 2})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if b == nil {
		t.Fatal("expected non-nil Blanker")
	}
}

func TestBlanker_Apply_EmptyInput(t *testing.T) {
	b, _ := New(Options{Every: 2})
	result := b.Apply([]string{})
	if len(result) != 0 {
		t.Fatalf("expected empty output, got %d lines", len(result))
	}
}

func TestBlanker_Apply_InsertsBlankLines(t *testing.T) {
	b, _ := New(Options{Every: 2})
	input := []string{"a", "b", "c", "d"}
	result := b.Apply(input)
	// Expect blank after line 2 and line 4
	if len(result) != 6 {
		t.Fatalf("expected 6 lines, got %d: %v", len(result), result)
	}
	if result[2] != "" {
		t.Errorf("expected blank at index 2, got %q", result[2])
	}
	if result[5] != "" {
		t.Errorf("expected blank at index 5, got %q", result[5])
	}
}

func TestBlanker_Apply_WithOffset(t *testing.T) {
	b, _ := New(Options{Every: 2, Offset: 2})
	input := []string{"a", "b", "c", "d", "e", "f"}
	result := b.Apply(input)
	// First 2 lines skipped for counting; blank inserted after every 2 from line index 2 onward
	// lines: a b c d [blank] e f [blank]
	for _, line := range result {
		_ = line
	}
	blankCount := 0
	for _, line := range result {
		if line == "" {
			blankCount++
		}
	}
	if blankCount != 2 {
		t.Errorf("expected 2 blank lines, got %d: %v", blankCount, result)
	}
}

func TestBlanker_Apply_OriginalUnmodified(t *testing.T) {
	b, _ := New(Options{Every: 2})
	input := []string{"x", "y", "z"}
	copy := append([]string{}, input...)
	b.Apply(input)
	for i, v := range input {
		if v != copy[i] {
			t.Errorf("original slice modified at index %d", i)
		}
	}
}

func TestFormatSummary(t *testing.T) {
	b, _ := New(Options{Every: 3, Offset: 1})
	s := FormatSummary(b, 9, 3)
	if !strings.Contains(s, "3 blank lines") {
		t.Errorf("unexpected summary: %s", s)
	}
	if !strings.Contains(s, "every=3") {
		t.Errorf("expected every=3 in summary: %s", s)
	}
}
