package freq

import (
	"testing"
)

func TestNew_NegativeTopN(t *testing.T) {
	_, err := New(Options{TopN: -1})
	if err == nil {
		t.Fatal("expected error for negative TopN")
	}
}

func TestNew_NegativeField(t *testing.T) {
	_, err := New(Options{Field: -1})
	if err == nil {
		t.Fatal("expected error for negative Field")
	}
}

func TestNew_ValidOptions(t *testing.T) {
	c, err := New(Options{TopN: 5, CaseInsensitive: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c == nil {
		t.Fatal("expected non-nil Counter")
	}
}

func TestCounter_Add_And_Total(t *testing.T) {
	c, _ := New(Options{})
	c.Add([]string{"foo", "bar", "foo", "baz", "foo"})
	if got := c.Total(); got != 5 {
		t.Fatalf("Total: want 5, got %d", got)
	}
}

func TestCounter_Results_Sorted(t *testing.T) {
	c, _ := New(Options{})
	c.Add([]string{"a", "b", "a", "c", "b", "a"})
	res := c.Results()
	if len(res) != 3 {
		t.Fatalf("expected 3 results, got %d", len(res))
	}
	if res[0].Value != "a" || res[0].Count != 3 {
		t.Errorf("first entry: want a/3, got %s/%d", res[0].Value, res[0].Count)
	}
}

func TestCounter_TopN_Limits(t *testing.T) {
	c, _ := New(Options{TopN: 2})
	c.Add([]string{"x", "y", "z", "x", "y", "x"})
	res := c.Results()
	if len(res) != 2 {
		t.Fatalf("expected 2 results with TopN=2, got %d", len(res))
	}
}

func TestCounter_CaseInsensitive(t *testing.T) {
	c, _ := New(Options{CaseInsensitive: true})
	c.Add([]string{"Hello", "hello", "HELLO"})
	res := c.Results()
	if len(res) != 1 {
		t.Fatalf("expected 1 unique key, got %d", len(res))
	}
	if res[0].Count != 3 {
		t.Errorf("expected count 3, got %d", res[0].Count)
	}
}

func TestCounter_FieldExtraction(t *testing.T) {
	c, _ := New(Options{Field: 2})
	c.Add([]string{"2024-01-01 ERROR msg", "2024-01-02 ERROR msg", "2024-01-03 INFO msg"})
	res := c.Results()
	if res[0].Value != "ERROR" || res[0].Count != 2 {
		t.Errorf("field extraction: want ERROR/2, got %s/%d", res[0].Value, res[0].Count)
	}
}

func TestCounter_FieldOutOfRange(t *testing.T) {
	c, _ := New(Options{Field: 10})
	c.Add([]string{"only three words here", "also short"})
	res := c.Results()
	// both lines produce empty key — one unique entry
	if len(res) != 1 {
		t.Fatalf("expected 1 bucket for out-of-range field, got %d", len(res))
	}
}

func TestCounter_EmptyInput(t *testing.T) {
	c, _ := New(Options{})
	c.Add([]string{})
	if c.Total() != 0 {
		t.Errorf("expected Total 0 for empty input")
	}
	if len(c.Results()) != 0 {
		t.Errorf("expected empty results for empty input")
	}
}
