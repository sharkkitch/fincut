package timestamp

import (
	"strings"
	"testing"
	"time"
)

func TestNew_MutuallyExclusive(t *testing.T) {
	_, err := New(Options{Prepend: true, Append: true})
	if err == nil {
		t.Fatal("expected error for mutually exclusive flags")
	}
}

func TestNew_NeitherSet(t *testing.T) {
	_, err := New(Options{})
	if err == nil {
		t.Fatal("expected error when neither Prepend nor Append is set")
	}
}

func TestNew_ValidPrepend(t *testing.T) {
	s, err := New(Options{Prepend: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s == nil {
		t.Fatal("expected non-nil Stamper")
	}
}

func TestNew_DefaultFormat(t *testing.T) {
	s, err := New(Options{Prepend: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.opts.Format != time.RFC3339 {
		t.Errorf("expected default format %q, got %q", time.RFC3339, s.opts.Format)
	}
}

func TestStamper_Apply_Prepend(t *testing.T) {
	s, _ := New(Options{Prepend: true, Format: "2006", Sep: "|"})
	year := time.Now().Format("2006")
	out := s.Apply([]string{"hello", "world"})
	for _, line := range out {
		if !strings.HasPrefix(line, year+"|") {
			t.Errorf("expected line to start with %q, got %q", year+"|", line)
		}
	}
}

func TestStamper_Apply_Append(t *testing.T) {
	s, _ := New(Options{Append: true, Format: "2006", Sep: "|"})
	year := time.Now().Format("2006")
	out := s.Apply([]string{"hello"})
	if !strings.HasSuffix(out[0], "|"+year) {
		t.Errorf("expected line to end with %q, got %q", "|"+year, out[0])
	}
}

func TestStamper_Apply_EmptyInput(t *testing.T) {
	s, _ := New(Options{Prepend: true})
	out := s.Apply([]string{})
	if len(out) != 0 {
		t.Errorf("expected empty output, got %d lines", len(out))
	}
}

func TestFormatSummary(t *testing.T) {
	summary := FormatSummary(42, time.RFC3339)
	if !strings.Contains(summary, "42") {
		t.Errorf("expected summary to contain line count, got %q", summary)
	}
}
