package redact

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

func TestNew_DefaultReplacement(t *testing.T) {
	r, err := New(Options{Patterns: []string{"secret"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.opts.Replacement != "[REDACTED]" {
		t.Errorf("expected default replacement, got %q", r.opts.Replacement)
	}
}

func TestNew_ValidOptions(t *testing.T) {
	r, err := New(Options{
		Patterns:    []string{`\d{4}-\d{4}-\d{4}-\d{4}`},
		Replacement: "***",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r == nil {
		t.Fatal("expected non-nil Redacter")
	}
}

func TestRedacter_Apply_WholeLineMatch(t *testing.T) {
	r, _ := New(Options{Patterns: []string{`password=.*`}})
	input := []string{"username=alice", "password=hunter2", "role=admin"}
	out := r.Apply(input)
	if out[0] != "username=alice" {
		t.Errorf("line 0 should be unchanged, got %q", out[0])
	}
	if out[1] != "[REDACTED]" {
		t.Errorf("line 1 should be redacted, got %q", out[1])
	}
	if out[2] != "role=admin" {
		t.Errorf("line 2 should be unchanged, got %q", out[2])
	}
}

func TestRedacter_Apply_PartialMatch(t *testing.T) {
	r, _ := New(Options{
		Patterns:     []string{`token`},
		MatchPartial: true,
	})
	input := []string{"bearer token abc123", "safe line"}
	out := r.Apply(input)
	if out[0] != "[REDACTED]" {
		t.Errorf("expected redaction, got %q", out[0])
	}
	if out[1] != "safe line" {
		t.Errorf("expected unchanged, got %q", out[1])
	}
}

func TestRedacter_Apply_CustomReplacement(t *testing.T) {
	r, _ := New(Options{
		Patterns:     []string{`secret`},
		Replacement:  "<hidden>",
		MatchPartial: true,
	})
	out := r.Apply([]string{"my secret value"})
	if out[0] != "<hidden>" {
		t.Errorf("expected <hidden>, got %q", out[0])
	}
}

func TestRedacter_Stats(t *testing.T) {
	r, _ := New(Options{
		Patterns:     []string{`DROP TABLE`},
		MatchPartial: true,
	})
	lines := []string{
		"SELECT * FROM users",
		"DROP TABLE users",
		"INSERT INTO logs VALUES(1)",
		"DROP TABLE sessions",
	}
	r.Apply(lines)
	total, redacted := r.Stats()
	if total != 4 {
		t.Errorf("expected total=4, got %d", total)
	}
	if redacted != 2 {
		t.Errorf("expected redacted=2, got %d", redacted)
	}
}

func TestRedacter_Apply_NoMatch(t *testing.T) {
	r, _ := New(Options{Patterns: []string{`nomatch`}, MatchPartial: true})
	input := []string{"line one", "line two"}
	out := r.Apply(input)
	for i, line := range out {
		if !strings.EqualFold(line, input[i]) {
			t.Errorf("line %d should be unchanged", i)
		}
	}
	_, redacted := r.Stats()
	if redacted != 0 {
		t.Errorf("expected 0 redactions, got %d", redacted)
	}
}
