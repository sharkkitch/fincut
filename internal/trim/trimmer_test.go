package trim_test

import (
	"strings"
	"testing"

	"github.com/user/fincut/internal/trim"
)

func TestNewTrimmer_InvalidOptions(t *testing.T) {
	_, err := trim.NewTrimmer(trim.Options{MaxLines: -1})
	if err == nil {
		t.Fatal("expected error for negative MaxLines")
	}
	_, err = trim.NewTrimmer(trim.Options{MaxBytes: -10})
	if err == nil {
		t.Fatal("expected error for negative MaxBytes")
	}
}

func TestNewTrimmer_ValidOptions(t *testing.T) {
	_, err := trim.NewTrimmer(trim.Options{MaxLines: 10, MaxBytes: 1024})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestTrimmer_Apply_MaxLines(t *testing.T) {
	tr, _ := trim.NewTrimmer(trim.Options{MaxLines: 2})
	lines := []string{"line1", "line2", "line3", "line4"}
	var buf strings.Builder
	_, err := tr.Apply(lines, &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := strings.Split(strings.TrimRight(buf.String(), "\n"), "\n")
	if len(got) != 2 {
		t.Errorf("expected 2 lines, got %d", len(got))
	}
}

func TestTrimmer_Apply_MaxBytes(t *testing.T) {
	tr, _ := trim.NewTrimmer(trim.Options{MaxBytes: 12})
	// each line is 5 chars + newline = 6 bytes; 2 lines = 12 bytes
	lines := []string{"hello", "world", "extra"}
	var buf strings.Builder
	tr.Apply(lines, &buf)
	got := strings.Split(strings.TrimRight(buf.String(), "\n"), "\n")
	if len(got) != 2 {
		t.Errorf("expected 2 lines within byte limit, got %d: %v", len(got), got)
	}
}

func TestTrimmer_Apply_StripEmpty(t *testing.T) {
	tr, _ := trim.NewTrimmer(trim.Options{StripEmpty: true})
	lines := []string{"line1", "", "  ", "line2"}
	var buf strings.Builder
	tr.Apply(lines, &buf)
	got := strings.Split(strings.TrimRight(buf.String(), "\n"), "\n")
	if len(got) != 2 {
		t.Errorf("expected 2 non-empty lines, got %d: %v", len(got), got)
	}
}

func TestTrimmer_Apply_NoLimits(t *testing.T) {
	tr, _ := trim.NewTrimmer(trim.Options{})
	lines := []string{"a", "b", "c"}
	var buf strings.Builder
	tr.Apply(lines, &buf)
	got := strings.Split(strings.TrimRight(buf.String(), "\n"), "\n")
	if len(got) != 3 {
		t.Errorf("expected all 3 lines, got %d", len(got))
	}
}

func TestTrimmer_Count(t *testing.T) {
	tr, _ := trim.NewTrimmer(trim.Options{MaxLines: 3})
	lines := []string{"a", "b", "c", "d", "e"}
	count := tr.Count(lines)
	if count != 3 {
		t.Errorf("expected count 3, got %d", count)
	}
}
