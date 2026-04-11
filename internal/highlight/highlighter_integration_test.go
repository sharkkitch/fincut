package highlight_test

import (
	"strings"
	"testing"

	"github.com/user/fincut/internal/highlight"
)

func TestHighlighter_MultiPattern_IndependentColors(t *testing.T) {
	h, err := highlight.New(highlight.Options{
		Patterns: []string{"error", "warn", "info"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	line := "error warn info"
	result := h.Apply(line)

	// All three terms should be highlighted
	stripped := highlight.StripANSI(result)
	for _, term := range []string{"error", "warn", "info"} {
		if !strings.Contains(stripped, term) {
			t.Errorf("expected term %q preserved after strip", term)
		}
	}
	// Should contain at least two distinct color codes
	if strings.Count(result, "\033[") < 2 {
		t.Errorf("expected multiple ANSI codes, got: %q", result)
	}
}

func TestHighlighter_RoundTrip_StripRestoresText(t *testing.T) {
	original := "2024-01-15 ERROR disk full WARN retry limit"
	h, _ := highlight.New(highlight.Options{
		Patterns: []string{"ERROR", "WARN"},
		Bold:     true,
	})
	highlighted := h.Apply(original)
	stripped := highlight.StripANSI(highlighted)
	if stripped != original {
		t.Errorf("round-trip failed:\n  want: %q\n   got: %q", original, stripped)
	}
}
