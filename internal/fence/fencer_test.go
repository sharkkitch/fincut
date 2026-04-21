package fence

import (
	"testing"
)

func TestNew_EmptyOpenPattern(t *testing.T) {
	_, err := New(Options{OpenPattern: "", ClosePattern: `---`})
	if err == nil {
		t.Fatal("expected error for empty open pattern")
	}
}

func TestNew_EmptyClosePattern(t *testing.T) {
	_, err := New(Options{OpenPattern: `---`, ClosePattern: ""})
	if err == nil {
		t.Fatal("expected error for empty close pattern")
	}
}

func TestNew_InvalidOpenPattern(t *testing.T) {
	_, err := New(Options{OpenPattern: `[bad`, ClosePattern: `---`})
	if err == nil {
		t.Fatal("expected error for invalid open pattern")
	}
}

func TestNew_ValidOptions(t *testing.T) {
	f, err := New(Options{OpenPattern: `^BEGIN`, ClosePattern: `^END`})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f == nil {
		t.Fatal("expected non-nil Fencer")
	}
}

func TestFencer_Apply_BasicExtract(t *testing.T) {
	f, _ := New(Options{OpenPattern: `^BEGIN`, ClosePattern: `^END`})
	input := []string{"preamble", "BEGIN", "inside1", "inside2", "END", "postamble"}
	got := f.Apply(input)
	if len(got) != 2 || got[0] != "inside1" || got[1] != "inside2" {
		t.Fatalf("unexpected result: %v", got)
	}
}

func TestFencer_Apply_IncludeDelim(t *testing.T) {
	f, _ := New(Options{OpenPattern: `^BEGIN`, ClosePattern: `^END`, IncludeDelim: true})
	input := []string{"BEGIN", "body", "END"}
	got := f.Apply(input)
	if len(got) != 3 {
		t.Fatalf("expected 3 lines, got %d: %v", len(got), got)
	}
}

func TestFencer_Apply_MultipleBlocks(t *testing.T) {
	f, _ := New(Options{OpenPattern: `^\[`, ClosePattern: `^\]`})
	input := []string{"[", "a", "]", "gap", "[", "b", "c", "]"}
	got := f.Apply(input)
	if len(got) != 3 || got[0] != "a" || got[1] != "b" || got[2] != "c" {
		t.Fatalf("unexpected result: %v", got)
	}
}

func TestFencer_Apply_EmptyInput(t *testing.T) {
	f, _ := New(Options{OpenPattern: `^BEGIN`, ClosePattern: `^END`})
	got := f.Apply([]string{})
	if len(got) != 0 {
		t.Fatalf("expected empty result, got %v", got)
	}
}

func TestFencer_Apply_UnclosedBlock(t *testing.T) {
	f, _ := New(Options{OpenPattern: `^BEGIN`, ClosePattern: `^END`})
	input := []string{"BEGIN", "orphan"}
	got := f.Apply(input)
	if len(got) != 1 || got[0] != "orphan" {
		t.Fatalf("unexpected result: %v", got)
	}
}

func TestFormatSummary(t *testing.T) {
	s := FormatSummary(100, 42)
	if s == "" {
		t.Fatal("expected non-empty summary")
	}
}
