package truncate

import (
	"strings"
	"testing"
)

func TestNewTruncator_NegativeMaxRunes(t *testing.T) {
	_, err := NewTruncator(Options{MaxRunes: -1})
	if err == nil {
		t.Fatal("expected error for negative MaxRunes, got nil")
	}
}

func TestNewTruncator_ZeroMaxRunes(t *testing.T) {
	tr, err := NewTruncator(Options{MaxRunes: 0})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tr == nil {
		t.Fatal("expected non-nil Truncator")
	}
}

func TestNewTruncator_DefaultEllipsis(t *testing.T) {
	tr, err := NewTruncator(Options{MaxRunes: 10})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tr.opts.Ellipsis != "..." {
		t.Errorf("expected default ellipsis '...', got %q", tr.opts.Ellipsis)
	}
}

func TestTruncator_Apply_NoLimit(t *testing.T) {
	tr, _ := NewTruncator(Options{MaxRunes: 0})
	lines := []string{"short", strings.Repeat("x", 200)}
	out := tr.Apply(lines)
	for i, line := range out {
		if line != lines[i] {
			t.Errorf("line %d: expected unchanged, got %q", i, line)
		}
	}
}

func TestTruncator_Apply_WithinLimit(t *testing.T) {
	tr, _ := NewTruncator(Options{MaxRunes: 20})
	lines := []string{"hello", "world"}
	out := tr.Apply(lines)
	for i, line := range out {
		if line != lines[i] {
			t.Errorf("line %d: expected %q, got %q", i, lines[i], line)
		}
	}
}

func TestTruncator_Apply_ExceedsLimit(t *testing.T) {
	tr, _ := NewTruncator(Options{MaxRunes: 10, Ellipsis: "..."})
	lines := []string{strings.Repeat("a", 20)}
	out := tr.Apply(lines)
	if len(out) != 1 {
		t.Fatalf("expected 1 line, got %d", len(out))
	}
	if len([]rune(out[0])) != 10 {
		t.Errorf("expected 10 runes, got %d: %q", len([]rune(out[0])), out[0])
	}
	if !strings.HasSuffix(out[0], "...") {
		t.Errorf("expected ellipsis suffix, got %q", out[0])
	}
}

func TestTruncator_Apply_CustomEllipsis(t *testing.T) {
	tr, _ := NewTruncator(Options{MaxRunes: 8, Ellipsis: "→"})
	lines := []string{"abcdefghij"}
	out := tr.Apply(lines)
	if !strings.HasSuffix(out[0], "→") {
		t.Errorf("expected custom ellipsis, got %q", out[0])
	}
	if len([]rune(out[0])) != 8 {
		t.Errorf("expected 8 runes, got %d", len([]rune(out[0])))
	}
}

func TestTruncator_Apply_MultibyteRunes(t *testing.T) {
	tr, _ := NewTruncator(Options{MaxRunes: 5, Ellipsis: "..."})
	// Each '日' is 3 bytes but 1 rune.
	lines := []string{"日本語テスト入力"}
	out := tr.Apply(lines)
	if len([]rune(out[0])) != 5 {
		t.Errorf("expected 5 runes, got %d: %q", len([]rune(out[0])), out[0])
	}
}
