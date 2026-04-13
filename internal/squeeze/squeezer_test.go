package squeeze

import (
	"strings"
	"testing"
)

func TestNew_DefaultSeparator(t *testing.T) {
	s, err := New(Options{Annotate: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.opts.Separator != "\t" {
		t.Errorf("expected tab separator, got %q", s.opts.Separator)
	}
}

func TestNew_CustomSeparator(t *testing.T) {
	s, err := New(Options{Annotate: true, Separator: " | "})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.opts.Separator != " | " {
		t.Errorf("expected ' | ', got %q", s.opts.Separator)
	}
}

func TestSqueezer_Apply_EmptyInput(t *testing.T) {
	s, _ := New(Options{})
	out := s.Apply([]string{})
	if len(out) != 0 {
		t.Errorf("expected empty output, got %d lines", len(out))
	}
}

func TestSqueezer_Apply_NoDuplicates(t *testing.T) {
	s, _ := New(Options{})
	input := []string{"alpha", "beta", "gamma"}
	out := s.Apply(input)
	if len(out) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(out))
	}
}

func TestSqueezer_Apply_AllSame(t *testing.T) {
	s, _ := New(Options{})
	input := []string{"same", "same", "same", "same"}
	out := s.Apply(input)
	if len(out) != 1 {
		t.Fatalf("expected 1 line, got %d", len(out))
	}
	if out[0] != "same" {
		t.Errorf("expected 'same', got %q", out[0])
	}
}

func TestSqueezer_Apply_MixedRuns(t *testing.T) {
	s, _ := New(Options{})
	input := []string{"a", "a", "b", "c", "c", "c", "a"}
	out := s.Apply(input)
	expected := []string{"a", "b", "c", "a"}
	if len(out) != len(expected) {
		t.Fatalf("expected %d lines, got %d", len(expected), len(out))
	}
	for i, e := range expected {
		if out[i] != e {
			t.Errorf("line %d: expected %q, got %q", i, e, out[i])
		}
	}
}

func TestSqueezer_Apply_Annotate(t *testing.T) {
	s, _ := New(Options{Annotate: true, Separator: ":"})
	input := []string{"x", "x", "x", "y", "y"}
	out := s.Apply(input)
	if len(out) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(out))
	}
	if !strings.HasPrefix(out[0], "3:") {
		t.Errorf("expected count prefix '3:', got %q", out[0])
	}
	if !strings.HasPrefix(out[1], "2:") {
		t.Errorf("expected count prefix '2:', got %q", out[1])
	}
}

func TestFormatSummary(t *testing.T) {
	summary := FormatSummary(10, 6)
	if !strings.Contains(summary, "4 removed") {
		t.Errorf("expected '4 removed' in summary, got %q", summary)
	}
}
