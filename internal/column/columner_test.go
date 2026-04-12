package column

import (
	"testing"
)

func TestNew_EmptyDelimiter(t *testing.T) {
	_, err := New(Options{Delimiter: "", Fields: []int{1}})
	if err == nil {
		t.Fatal("expected error for empty delimiter")
	}
}

func TestNew_NoFields(t *testing.T) {
	_, err := New(Options{Delimiter: ",", Fields: []int{}})
	if err == nil {
		t.Fatal("expected error for empty fields")
	}
}

func TestNew_ZeroFieldIndex(t *testing.T) {
	_, err := New(Options{Delimiter: ",", Fields: []int{0}})
	if err == nil {
		t.Fatal("expected error for field index 0")
	}
}

func TestNew_ValidOptions(t *testing.T) {
	c, err := New(Options{Delimiter: ",", Fields: []int{1, 3}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c == nil {
		t.Fatal("expected non-nil Columner")
	}
}

func TestColumner_Apply_BasicExtract(t *testing.T) {
	c, _ := New(Options{Delimiter: ",", Fields: []int{1, 3}})
	lines := []string{"a,b,c", "x,y,z"}
	got, err := c.Apply(lines)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("expected 2 results, got %d", len(got))
	}
	if got[0] != "a,c" {
		t.Errorf("expected \"a,c\", got %q", got[0])
	}
	if got[1] != "x,z" {
		t.Errorf("expected \"x,z\", got %q", got[1])
	}
}

func TestColumner_Apply_CustomOutputSep(t *testing.T) {
	c, _ := New(Options{Delimiter: ",", Fields: []int{1, 2}, OutputSep: "|" })
	got, err := c.Apply([]string{"foo,bar,baz"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got[0] != "foo|bar" {
		t.Errorf("expected \"foo|bar\", got %q", got[0])
	}
}

func TestColumner_Apply_OutOfRange_Lenient(t *testing.T) {
	c, _ := New(Options{Delimiter: ",", Fields: []int{1, 5}})
	got, err := c.Apply([]string{"a,b,c"})
	if err != nil {
		t.Fatalf("unexpected error in lenient mode: %v", err)
	}
	if got[0] != "a," {
		t.Errorf("expected \"a,\", got %q", got[0])
	}
}

func TestColumner_Apply_OutOfRange_Strict(t *testing.T) {
	c, _ := New(Options{Delimiter: ",", Fields: []int{1, 5}, Strict: true})
	_, err := c.Apply([]string{"a,b,c"})
	if err == nil {
		t.Fatal("expected error in strict mode for missing field")
	}
}

func TestColumner_Apply_EmptyInput(t *testing.T) {
	c, _ := New(Options{Delimiter: "\t", Fields: []int{1}})
	got, err := c.Apply([]string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 0 {
		t.Errorf("expected empty result, got %d lines", len(got))
	}
}
