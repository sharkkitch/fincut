package extract

import (
	"strings"
	"testing"
)

func TestNew_EmptyPattern(t *testing.T) {
	_, err := New(Options{Pattern: ""})
	if err == nil {
		t.Fatal("expected error for empty pattern")
	}
}

func TestNew_InvalidPattern(t *testing.T) {
	_, err := New(Options{Pattern: "([invalid"})
	if err == nil {
		t.Fatal("expected error for invalid regex")
	}
}

func TestNew_NoCaptureGroup(t *testing.T) {
	_, err := New(Options{Pattern: "nogroup"})
	if err == nil {
		t.Fatal("expected error when no capture group present")
	}
}

func TestNew_GroupExceedsCaptures(t *testing.T) {
	_, err := New(Options{Pattern: `(\w+)`, Group: 5})
	if err == nil {
		t.Fatal("expected error when group exceeds available captures")
	}
}

func TestNew_ValidOptions(t *testing.T) {
	e, err := New(Options{Pattern: `level=(\w+)`})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if e == nil {
		t.Fatal("expected non-nil extractor")
	}
}

func TestNew_DefaultGroupIsOne(t *testing.T) {
	e, err := New(Options{Pattern: `(\w+)`})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if e.group != 1 {
		t.Fatalf("expected group 1, got %d", e.group)
	}
}

func TestExtractor_Apply_BasicExtract(t *testing.T) {
	e, _ := New(Options{Pattern: `level=(\w+)`, SkipUnmatched: true})
	lines := []string{"level=info msg=hello", "level=error msg=fail", "no match here"}
	got := e.Apply(lines)
	if len(got) != 2 {
		t.Fatalf("expected 2 results, got %d", len(got))
	}
	if got[0] != "info" || got[1] != "error" {
		t.Fatalf("unexpected results: %v", got)
	}
}

func TestExtractor_Apply_PassThroughUnmatched(t *testing.T) {
	e, _ := New(Options{Pattern: `level=(\w+)`, SkipUnmatched: false})
	lines := []string{"level=info msg=hello", "no match here"}
	got := e.Apply(lines)
	if len(got) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(got))
	}
	if got[1] != "no match here" {
		t.Fatalf("expected unmatched line passed through, got %q", got[1])
	}
}

func TestExtractor_Apply_EmptyInput(t *testing.T) {
	e, _ := New(Options{Pattern: `(\w+)`})
	got := e.Apply([]string{})
	if len(got) != 0 {
		t.Fatalf("expected empty output, got %d lines", len(got))
	}
}

func TestExtractor_Apply_SecondGroup(t *testing.T) {
	e, _ := New(Options{Pattern: `(\w+)=(\w+)`, Group: 2})
	lines := []string{"level=info", "host=localhost"}
	got := e.Apply(lines)
	for _, g := range got {
		if strings.Contains(g, "=") {
			t.Fatalf("expected extracted value only, got %q", g)
		}
	}
}

func TestFormatSummary(t *testing.T) {
	s := FormatSummary(10, 7, 3)
	if !strings.Contains(s, "7/10") {
		t.Fatalf("expected summary to contain 7/10, got %q", s)
	}
}
