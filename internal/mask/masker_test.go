package mask

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

func TestNew_InvalidPattern(t *testing.T) {
	_, err := New(Options{Patterns: []string{"[invalid"}})
	if err == nil {
		t.Fatal("expected error for invalid regex")
	}
}

func TestNew_ValidOptions(t *testing.T) {
	m, err := New(Options{Patterns: []string{`\d+`}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m == nil {
		t.Fatal("expected non-nil masker")
	}
}

func TestNew_DefaultReplacement(t *testing.T) {
	m, err := New(Options{Patterns: []string{`\d+`}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m.replacement != "[REDACTED]" {
		t.Errorf("expected default replacement, got %q", m.replacement)
	}
}

func TestMasker_Apply_RedactsMatch(t *testing.T) {
	m, _ := New(Options{Patterns: []string{`\d{4}-\d{2}-\d{2}`}})
	lines := []string{"event on 2024-01-15 occurred", "no date here"}
	out := m.Apply(lines)
	if strings.Contains(out[0], "2024-01-15") {
		t.Errorf("expected date to be redacted, got %q", out[0])
	}
	if out[1] != "no date here" {
		t.Errorf("expected unchanged line, got %q", out[1])
	}
}

func TestMasker_Apply_CustomReplacement(t *testing.T) {
	m, _ := New(Options{Patterns: []string{`password=\S+`}, Replacement: "***"})
	out := m.Apply([]string{"login password=secret123"})
	if !strings.Contains(out[0], "***") {
		t.Errorf("expected *** in output, got %q", out[0])
	}
}

func TestMasker_Apply_MultiplePatterns(t *testing.T) {
	m, _ := New(Options{
		Patterns: []string{`\b\d{3}-\d{2}-\d{4}\b`, `token=[A-Za-z0-9]+`},
	})
	line := "ssn=123-45-6789 token=abc123XYZ"
	out := m.Apply([]string{line})
	if strings.Contains(out[0], "123-45-6789") {
		t.Errorf("SSN not redacted: %q", out[0])
	}
	if strings.Contains(out[0], "abc123XYZ") {
		t.Errorf("token not redacted: %q", out[0])
	}
}

func TestMasker_CountRedacted(t *testing.T) {
	m, _ := New(Options{Patterns: []string{`secret`}})
	lines := []string{"this is secret", "nothing here", "another secret value"}
	count := m.CountRedacted(lines)
	if count != 2 {
		t.Errorf("expected 2 redacted lines, got %d", count)
	}
}
