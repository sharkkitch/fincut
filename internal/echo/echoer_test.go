package echo

import (
	"bytes"
	"strings"
	"testing"
)

func TestNew_NilWriter(t *testing.T) {
	_, err := New(Options{Writer: nil})
	if err == nil {
		t.Fatal("expected error for nil Writer, got nil")
	}
}

func TestNew_ValidOptions(t *testing.T) {
	var buf bytes.Buffer
	e, err := New(Options{Writer: &buf})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if e == nil {
		t.Fatal("expected non-nil Echoer")
	}
}

func TestEchoer_Apply_PassThrough(t *testing.T) {
	var buf bytes.Buffer
	e, _ := New(Options{Writer: &buf})

	input := []string{"alpha", "beta", "gamma"}
	out := e.Apply(input)

	if len(out) != len(input) {
		t.Fatalf("expected %d lines, got %d", len(input), len(out))
	}
	for i, line := range input {
		if out[i] != line {
			t.Errorf("line %d: expected %q, got %q", i, line, out[i])
		}
	}
}

func TestEchoer_Apply_WritesToSecondary(t *testing.T) {
	var buf bytes.Buffer
	e, _ := New(Options{Writer: &buf})

	e.Apply([]string{"hello", "world"})

	got := buf.String()
	if !strings.Contains(got, "hello") || !strings.Contains(got, "world") {
		t.Errorf("expected both lines in secondary output, got: %q", got)
	}
}

func TestEchoer_Apply_WithPrefix(t *testing.T) {
	var buf bytes.Buffer
	e, _ := New(Options{Writer: &buf, Prefix: "[DBG] "})

	e.Apply([]string{"trace line"})

	if !strings.Contains(buf.String(), "[DBG] trace line") {
		t.Errorf("expected prefix in output, got: %q", buf.String())
	}
}

func TestEchoer_Apply_EmptyInput(t *testing.T) {
	var buf bytes.Buffer
	e, _ := New(Options{Writer: &buf})

	out := e.Apply([]string{})
	if len(out) != 0 {
		t.Errorf("expected empty output, got %d lines", len(out))
	}
	if buf.Len() != 0 {
		t.Errorf("expected empty secondary output, got %q", buf.String())
	}
}

func TestEchoer_FormatSummary(t *testing.T) {
	var buf bytes.Buffer
	e, _ := New(Options{Writer: &buf})

	e.Apply([]string{"a", "b", "c"})

	summary := e.FormatSummary()
	if !strings.Contains(summary, "3") {
		t.Errorf("expected count 3 in summary, got: %q", summary)
	}
}
