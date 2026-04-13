package offset

import (
	"testing"
)

func TestNew_InvalidStartLine(t *testing.T) {
	_, err := New(Options{StartLine: 0})
	if err == nil {
		t.Fatal("expected error for StartLine=0")
	}
}

func TestNew_EndBeforeStart(t *testing.T) {
	_, err := New(Options{StartLine: 5, EndLine: 3})
	if err == nil {
		t.Fatal("expected error when EndLine < StartLine")
	}
}

func TestNew_ValidOptions(t *testing.T) {
	_, err := New(Options{StartLine: 1, EndLine: 10})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNew_ZeroEndLine(t *testing.T) {
	_, err := New(Options{StartLine: 1, EndLine: 0})
	if err != nil {
		t.Fatalf("EndLine=0 should be valid (no limit): %v", err)
	}
}

func TestOffsetter_Apply_AllLines(t *testing.T) {
	o, _ := New(Options{StartLine: 1})
	lines := []string{"alpha", "beta", "gamma"}
	entries := o.Apply(lines)
	if len(entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(entries))
	}
	if entries[0].ByteStart != 0 {
		t.Errorf("first entry ByteStart should be 0, got %d", entries[0].ByteStart)
	}
	if entries[1].ByteStart != int64(len("alpha")+1) {
		t.Errorf("second entry ByteStart wrong: %d", entries[1].ByteStart)
	}
}

func TestOffsetter_Apply_Range(t *testing.T) {
	o, _ := New(Options{StartLine: 2, EndLine: 3})
	lines := []string{"one", "two", "three", "four"}
	entries := o.Apply(lines)
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	if entries[0].Line != 2 || entries[0].Content != "two" {
		t.Errorf("unexpected first entry: %+v", entries[0])
	}
	if entries[1].Line != 3 || entries[1].Content != "three" {
		t.Errorf("unexpected second entry: %+v", entries[1])
	}
}

func TestOffsetter_Apply_EmptyInput(t *testing.T) {
	o, _ := New(Options{StartLine: 1})
	entries := o.Apply([]string{})
	if len(entries) != 0 {
		t.Fatalf("expected 0 entries for empty input, got %d", len(entries))
	}
}

func TestOffsetter_Apply_ByteOffsets_Consistent(t *testing.T) {
	o, _ := New(Options{StartLine: 1})
	lines := []string{"hello", "world"}
	entries := o.Apply(lines)
	for _, e := range entries {
		expectedLen := int64(len(e.Content))
		gotLen := e.ByteEnd - e.ByteStart
		if gotLen != expectedLen {
			t.Errorf("line %d: expected byte span %d, got %d", e.Line, expectedLen, gotLen)
		}
	}
}
