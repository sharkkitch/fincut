package linefilter_test

import (
	"testing"

	"github.com/your-org/fincut/internal/linefilter"
)

func TestNew_NeitherSet(t *testing.T) {
	_, err := linefilter.New(linefilter.Options{})
	if err == nil {
		t.Fatal("expected error for empty options")
	}
}

func TestNew_InvalidIncludePattern(t *testing.T) {
	_, err := linefilter.New(linefilter.Options{Include: []string{"[invalid"}})
	if err == nil {
		t.Fatal("expected error for invalid include pattern")
	}
}

func TestNew_InvalidExcludePattern(t *testing.T) {
	_, err := linefilter.New(linefilter.Options{Exclude: []string{"[bad"}})
	if err == nil {
		t.Fatal("expected error for invalid exclude pattern")
	}
}

func TestNew_ValidOptions(t *testing.T) {
	_, err := linefilter.New(linefilter.Options{Include: []string{"foo"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestLineFilter_Apply_IncludeOnly(t *testing.T) {
	lf, _ := linefilter.New(linefilter.Options{Include: []string{"ERROR"}})
	input := []string{"INFO start", "ERROR boom", "DEBUG tick", "ERROR crash"}
	out := lf.Apply(input)
	if len(out) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(out))
	}
}

func TestLineFilter_Apply_ExcludeOnly(t *testing.T) {
	lf, _ := linefilter.New(linefilter.Options{Exclude: []string{"DEBUG"}})
	input := []string{"INFO start", "DEBUG tick", "ERROR crash"}
	out := lf.Apply(input)
	if len(out) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(out))
	}
}

func TestLineFilter_Apply_IncludeAndExclude(t *testing.T) {
	lf, _ := linefilter.New(linefilter.Options{
		Include: []string{"ERROR"},
		Exclude: []string{"transient"},
	})
	input := []string{"ERROR boom", "ERROR transient skip", "INFO ok"}
	out := lf.Apply(input)
	if len(out) != 1 || out[0] != "ERROR boom" {
		t.Fatalf("unexpected output: %v", out)
	}
}

func TestLineFilter_Apply_EmptyInput(t *testing.T) {
	lf, _ := linefilter.New(linefilter.Options{Include: []string{"x"}})
	out := lf.Apply(nil)
	if len(out) != 0 {
		t.Fatalf("expected empty output, got %v", out)
	}
}

func TestFormatSummary(t *testing.T) {
	s := linefilter.FormatSummary(10, 4)
	if s == "" {
		t.Fatal("expected non-empty summary")
	}
}
