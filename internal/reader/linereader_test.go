package reader_test

import (
	"strings"
	"testing"

	"github.com/yourorg/fincut/internal/reader"
)

func TestNewLineReader_InvalidOffsets(t *testing.T) {
	_, err := reader.NewLineReader(-1, -1)
	if err == nil {
		t.Fatal("expected error for negative start offset")
	}

	_, err = reader.NewLineReader(10, 5)
	if err == nil {
		t.Fatal("expected error when end < start")
	}
}

func TestNewLineReader_ValidOffsets(t *testing.T) {
	_, err := reader.NewLineReader(0, -1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, err = reader.NewLineReader(0, 100)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestLineReader_ReadFrom_AllLines(t *testing.T) {
	input := "line1\nline2\nline3\n"
	lr, _ := reader.NewLineReader(0, -1)

	lines, err := lr.ReadFrom(strings.NewReader(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(lines))
	}
	if lines[0] != "line1" || lines[2] != "line3" {
		t.Errorf("unexpected lines: %v", lines)
	}
}

func TestLineReader_ReadFrom_ByteRange(t *testing.T) {
	// "line1\n" is 6 bytes; start at 6 to skip first line
	input := "line1\nline2\nline3\n"
	lr, _ := reader.NewLineReader(6, -1)

	lines, err := lr.ReadFrom(strings.NewReader(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(lines) != 2 {
		t.Fatalf("expected 2 lines after skip, got %d", len(lines))
	}
	if lines[0] != "line2" {
		t.Errorf("expected first line to be 'line2', got %q", lines[0])
	}
}

func TestLineReader_ReadFrom_EndOffset(t *testing.T) {
	// Read only the first 6 bytes: "line1\n"
	input := "line1\nline2\nline3\n"
	lr, _ := reader.NewLineReader(0, 6)

	lines, err := lr.ReadFrom(strings.NewReader(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(lines) != 1 {
		t.Fatalf("expected 1 line, got %d", len(lines))
	}
	if lines[0] != "line1" {
		t.Errorf("expected 'line1', got %q", lines[0])
	}
}

func TestLineReader_ReadFrom_Empty(t *testing.T) {
	lr, _ := reader.NewLineReader(0, -1)
	lines, err := lr.ReadFrom(strings.NewReader(""))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(lines) != 0 {
		t.Fatalf("expected 0 lines, got %d", len(lines))
	}
}
