package transpose

import (
	"testing"
)

func TestNew_EmptyDelimiter(t *testing.T) {
	_, err := New(Options{Delimiter: ""})
	if err == nil {
		t.Fatal("expected error for empty delimiter")
	}
}

func TestNew_ValidOptions(t *testing.T) {
	tr, err := New(Options{Delimiter: ","})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tr == nil {
		t.Fatal("expected non-nil Transposer")
	}
}

func TestTransposer_Apply_EmptyInput(t *testing.T) {
	tr, _ := New(Options{Delimiter: ","})
	out := tr.Apply([]string{})
	if len(out) != 0 {
		t.Fatalf("expected empty output, got %v", out)
	}
}

func TestTransposer_Apply_SingleLine(t *testing.T) {
	tr, _ := New(Options{Delimiter: ","})
	out := tr.Apply([]string{"a,b,c"})
	// Each field becomes its own row.
	if len(out) != 3 {
		t.Fatalf("expected 3 output lines, got %d", len(out))
	}
	if out[0] != "a" || out[1] != "b" || out[2] != "c" {
		t.Fatalf("unexpected output: %v", out)
	}
}

func TestTransposer_Apply_Matrix(t *testing.T) {
	tr, _ := New(Options{Delimiter: ","})
	input := []string{"1,2,3", "4,5,6"}
	out := tr.Apply(input)
	// Expected: ["1,4", "2,5", "3,6"]
	if len(out) != 3 {
		t.Fatalf("expected 3 rows, got %d", len(out))
	}
	expected := []string{"1,4", "2,5", "3,6"}
	for i, want := range expected {
		if out[i] != want {
			t.Errorf("row %d: want %q, got %q", i, want, out[i])
		}
	}
}

func TestTransposer_Apply_JaggedNoPad(t *testing.T) {
	tr, _ := New(Options{Delimiter: ","})
	input := []string{"a,b,c", "x,y"}
	out := tr.Apply(input)
	if len(out) != 3 {
		t.Fatalf("expected 3 rows, got %d", len(out))
	}
	// Third row: only first line contributed a value.
	if out[2] != "c," {
		t.Errorf("row 2: want %q, got %q", "c,", out[2])
	}
}

func TestTransposer_Apply_JaggedWithPad(t *testing.T) {
	tr, _ := New(Options{Delimiter: ",", PadFields: true, FillEmpty: "N/A"})
	input := []string{"a,b,c", "x,y"}
	out := tr.Apply(input)
	if out[2] != "c,N/A" {
		t.Errorf("row 2: want %q, got %q", "c,N/A", out[2])
	}
}

func TestFormatSummary(t *testing.T) {
	s := FormatSummary(4, 3)
	if s == "" {
		t.Fatal("expected non-empty summary")
	}
}
