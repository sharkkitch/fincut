package index

import (
	"strings"
	"testing"
)

func TestBuild_NilReader(t *testing.T) {
	_, err := Build(nil)
	if err == nil {
		t.Fatal("expected error for nil reader, got nil")
	}
}

func TestBuild_EmptyInput(t *testing.T) {
	idx, err := Build(strings.NewReader(""))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if idx.Len() != 0 {
		t.Errorf("expected 0 entries, got %d", idx.Len())
	}
}

func TestBuild_LineCount(t *testing.T) {
	input := "alpha\nbeta\ngamma\n"
	idx, err := Build(strings.NewReader(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if idx.Len() != 3 {
		t.Errorf("expected 3 entries, got %d", idx.Len())
	}
}

func TestLookup_ValidLine(t *testing.T) {
	idx, _ := Build(strings.NewReader("hello\nworld\n"))
	e, err := idx.Lookup(0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if e.Offset != 0 {
		t.Errorf("expected offset 0, got %d", e.Offset)
	}
	if e.Length != 6 { // "hello" + newline
		t.Errorf("expected length 6, got %d", e.Length)
	}
}

func TestLookup_OutOfRange(t *testing.T) {
	idx, _ := Build(strings.NewReader("only\n"))
	_, err := idx.Lookup(5)
	if err == nil {
		t.Fatal("expected out-of-range error, got nil")
	}
}

func TestRange_ValidRange(t *testing.T) {
	idx, _ := Build(strings.NewReader("a\nb\nc\nd\n"))
	entries, err := idx.Range(1, 2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 2 {
		t.Errorf("expected 2 entries, got %d", len(entries))
	}
	if entries[0].Line != 1 {
		t.Errorf("expected first entry line 1, got %d", entries[0].Line)
	}
}

func TestRange_InvalidRange(t *testing.T) {
	idx, _ := Build(strings.NewReader("x\ny\n"))
	_, err := idx.Range(1, 0)
	if err == nil {
		t.Fatal("expected error for inverted range, got nil")
	}
}

func TestOffsets_Sequential(t *testing.T) {
	idx, _ := Build(strings.NewReader("ab\ncd\nef\n"))
	expected := []int64{0, 3, 6}
	for i, want := range expected {
		e, _ := idx.Lookup(i)
		if e.Offset != want {
			t.Errorf("line %d: expected offset %d, got %d", i, want, e.Offset)
		}
	}
}
