package annotate

import (
	"strings"
	"testing"
)

func TestNew_NoAnnotations(t *testing.T) {
	_, err := New(Options{})
	if err == nil {
		t.Fatal("expected error for empty options, got nil")
	}
}

func TestNew_ValidOptions(t *testing.T) {
	a, err := New(Options{LineNumbers: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if a == nil {
		t.Fatal("expected non-nil Annotator")
	}
}

func TestAnnotator_Apply_LineNumbers(t *testing.T) {
	a, _ := New(Options{LineNumbers: true})
	lines := []string{"alpha", "beta", "gamma"}
	out := a.Apply(lines)

	if len(out) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(out))
	}
	for i, line := range out {
		expected := strings.HasPrefix(line, strings.Join([]string{string(rune('1'+i))}, ""))
		if !expected {
			t.Errorf("line %d: expected line number prefix, got %q", i+1, line)
		}
	}
}

func TestAnnotator_Apply_Prefix(t *testing.T) {
	a, _ := New(Options{Prefix: "APP"})
	out := a.Apply([]string{"log line"})

	if !strings.Contains(out[0], "APP") {
		t.Errorf("expected prefix APP in output, got %q", out[0])
	}
	if !strings.HasSuffix(out[0], "log line") {
		t.Errorf("expected original line at end, got %q", out[0])
	}
}

func TestAnnotator_Apply_CustomSeparator(t *testing.T) {
	a, _ := New(Options{Prefix: "X", Separator: "::"})
	out := a.Apply([]string{"data"})

	if !strings.Contains(out[0], "::") {
		t.Errorf("expected '::' separator in output, got %q", out[0])
	}
}

func TestAnnotator_Apply_Timestamp(t *testing.T) {
	a, _ := New(Options{Timestamp: true})
	out := a.Apply([]string{"entry"})

	// RFC3339 timestamps contain 'T' and 'Z'
	if !strings.Contains(out[0], "T") {
		t.Errorf("expected RFC3339 timestamp in output, got %q", out[0])
	}
}

func TestAnnotator_Apply_EmptyInput(t *testing.T) {
	a, _ := New(Options{LineNumbers: true})
	out := a.Apply([]string{})

	if len(out) != 0 {
		t.Errorf("expected empty output, got %d lines", len(out))
	}
}

func TestAnnotator_Apply_CombinedOptions(t *testing.T) {
	a, _ := New(Options{LineNumbers: true, Prefix: "SVC", Separator: " | "})
	out := a.Apply([]string{"message"})

	if !strings.HasPrefix(out[0], "1") {
		t.Errorf("expected line number at start, got %q", out[0])
	}
	if !strings.Contains(out[0], "SVC") {
		t.Errorf("expected prefix SVC in output, got %q", out[0])
	}
	if !strings.HasSuffix(out[0], "message") {
		t.Errorf("expected original message at end, got %q", out[0])
	}
}
