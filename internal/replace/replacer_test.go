package replace

import (
	"strings"
	"testing"
)

func TestNew_NoPatterns(t *testing.T) {
	_, err := New(Options{})
	if err == nil {
		t.Fatal("expected error for empty patterns")
	}
}

func TestNew_MissingSeparator(t *testing.T) {
	_, err := New(Options{Patterns: []string{"noequals"}})
	if err == nil {
		t.Fatal("expected error for missing '=' separator")
	}
}

func TestNew_InvalidPattern(t *testing.T) {
	_, err := New(Options{Patterns: []string{"[invalid=(replacement)"}})
	if err == nil {
		t.Fatal("expected error for invalid regex")
	}
}

func TestNew_ValidOptions(t *testing.T) {
	r, err := New(Options{Patterns: []string{`\d+=NUM`}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r == nil {
		t.Fatal("expected non-nil Replacer")
	}
}

func TestReplacer_Apply_BasicSubstitution(t *testing.T) {
	r, _ := New(Options{Patterns: []string{"foo=bar"}})
	out := r.Apply([]string{"foo is foo", "no match here"})
	if out[0] != "bar is bar" {
		t.Errorf("got %q, want %q", out[0], "bar is bar")
	}
	if out[1] != "no match here" {
		t.Errorf("got %q, want %q", out[1], "no match here")
	}
}

func TestReplacer_Apply_RegexGroup(t *testing.T) {
	r, _ := New(Options{Patterns: []string{`(\w+)@example\.com=$1@replaced.org`}})
	out := r.Apply([]string{"contact user@example.com today"})
	if !strings.Contains(out[0], "user@replaced.org") {
		t.Errorf("got %q, expected backreference expansion", out[0])
	}
}

func TestReplacer_Apply_LiteralReplacement(t *testing.T) {
	r, _ := New(Options{
		Patterns: []string{`\d+=$1`},
		Literal:  true,
	})
	out := r.Apply([]string{"value 42 here"})
	// With literal mode, '$1' should appear verbatim, not be expanded.
	if !strings.Contains(out[0], "$1") {
		t.Errorf("got %q, expected literal '$1' in output", out[0])
	}
}

func TestReplacer_Apply_MultiplePatterns(t *testing.T) {
	r, _ := New(Options{Patterns: []string{"hello=hi", "world=earth"}})
	out := r.Apply([]string{"hello world"})
	if out[0] != "hi earth" {
		t.Errorf("got %q, want %q", out[0], "hi earth")
	}
}

func TestReplacer_Apply_EmptyInput(t *testing.T) {
	r, _ := New(Options{Patterns: []string{"a=b"}})
	out := r.Apply([]string{})
	if len(out) != 0 {
		t.Errorf("expected empty output, got %d lines", len(out))
	}
}
