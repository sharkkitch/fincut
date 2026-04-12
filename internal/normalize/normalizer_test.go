package normalize

import (
	"testing"
)

func TestNew_NoOptions(t *testing.T) {
	_, err := New(Options{})
	if err == nil {
		t.Fatal("expected error when no options set, got nil")
	}
}

func TestNew_ValidOptions(t *testing.T) {
	n, err := New(Options{TrimSpace: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n == nil {
		t.Fatal("expected non-nil Normalizer")
	}
}

func TestNormalizer_Apply_TrimSpace(t *testing.T) {
	n, _ := New(Options{TrimSpace: true})
	input := []string{"  hello  ", "\tworld\t", "no-trim"}
	got := n.Apply(input)
	want := []string{"hello", "world", "no-trim"}
	for i, g := range got {
		if g != want[i] {
			t.Errorf("line %d: got %q, want %q", i, g, want[i])
		}
	}
}

func TestNormalizer_Apply_Lowercase(t *testing.T) {
	n, _ := New(Options{Lowercase: true})
	input := []string{"ERROR: Disk Full", "WARN: High CPU", "info: ok"}
	got := n.Apply(input)
	want := []string{"error: disk full", "warn: high cpu", "info: ok"}
	for i, g := range got {
		if g != want[i] {
			t.Errorf("line %d: got %q, want %q", i, g, want[i])
		}
	}
}

func TestNormalizer_Apply_CollapseSpaces(t *testing.T) {
	n, _ := New(Options{CollapseSpaces: true})
	input := []string{"foo   bar", "a  b  c", "single"}
	got := n.Apply(input)
	want := []string{"foo bar", "a b c", "single"}
	for i, g := range got {
		if g != want[i] {
			t.Errorf("line %d: got %q, want %q", i, g, want[i])
		}
	}
}

func TestNormalizer_Apply_StripControl(t *testing.T) {
	n, _ := New(Options{StripControl: true})
	input := []string{"hello\x00world", "foo\x1bbar", "clean"}
	got := n.Apply(input)
	want := []string{"helloworld", "foobar", "clean"}
	for i, g := range got {
		if g != want[i] {
			t.Errorf("line %d: got %q, want %q", i, g, want[i])
		}
	}
}

func TestNormalizer_Apply_Combined(t *testing.T) {
	n, _ := New(Options{
		TrimSpace:      true,
		CollapseSpaces: true,
		Lowercase:      true,
	})
	input := []string{"  ERROR   Disk  Full  "}
	got := n.Apply(input)
	want := "error disk full"
	if got[0] != want {
		t.Errorf("got %q, want %q", got[0], want)
	}
}

func TestNormalizer_Summary(t *testing.T) {
	n, _ := New(Options{TrimSpace: true, Lowercase: true})
	s := n.Summary()
	if s == "" {
		t.Error("expected non-empty summary")
	}
}
