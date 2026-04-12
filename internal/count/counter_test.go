package count

import (
	"testing"
)

func TestNew_NegativeTopN(t *testing.T) {
	_, err := New(Options{TopN: -1})
	if err == nil {
		t.Fatal("expected error for negative TopN, got nil")
	}
}

func TestNew_ValidOptions(t *testing.T) {
	c, err := New(Options{TopN: 5})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c == nil {
		t.Fatal("expected non-nil Counter")
	}
}

func TestCounter_Add_And_Total(t *testing.T) {
	c, _ := New(Options{})
	c.Add("alpha")
	c.Add("beta")
	c.Add("alpha")
	if got := c.Total(); got != 3 {
		t.Errorf("Total() = %d, want 3", got)
	}
}

func TestCounter_Results_Sorted(t *testing.T) {
	c, _ := New(Options{})
	for i := 0; i < 3; i++ {
		c.Add("common")
	}
	c.Add("rare")

	results := c.Results()
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	if results[0].Line != "common" || results[0].Count != 3 {
		t.Errorf("first result = %+v, want {common 3}", results[0])
	}
	if results[1].Line != "rare" || results[1].Count != 1 {
		t.Errorf("second result = %+v, want {rare 1}", results[1])
	}
}

func TestCounter_TopN_Limits_Results(t *testing.T) {
	c, _ := New(Options{TopN: 2})
	lines := []string{"a", "b", "c", "a", "b", "a"}
	for _, l := range lines {
		c.Add(l)
	}
	results := c.Results()
	if len(results) != 2 {
		t.Errorf("expected 2 results with TopN=2, got %d", len(results))
	}
}

func TestCounter_CaseInsensitive(t *testing.T) {
	c, _ := New(Options{CaseInsensitive: true})
	c.Add("ERROR")
	c.Add("error")
	c.Add("Error")

	results := c.Results()
	if len(results) != 1 {
		t.Fatalf("expected 1 unique entry, got %d", len(results))
	}
	if results[0].Count != 3 {
		t.Errorf("Count = %d, want 3", results[0].Count)
	}
}

func TestCounter_EmptyInput(t *testing.T) {
	c, _ := New(Options{})
	if total := c.Total(); total != 0 {
		t.Errorf("Total() = %d, want 0", total)
	}
	if results := c.Results(); len(results) != 0 {
		t.Errorf("Results() length = %d, want 0", len(results))
	}
}

func TestCounter_Results_Sorted_StableOrder(t *testing.T) {
	// When two lines have equal counts, verify Results() returns a consistent
	// (deterministic) ordering across multiple calls.
	c, _ := New(Options{})
	c.Add("apple")
	c.Add("banana")

	first := c.Results()
	second := c.Results()

	if len(first) != len(second) {
		t.Fatalf("Results() length mismatch: %d vs %d", len(first), len(second))
	}
	for i := range first {
		if first[i].Line != second[i].Line || first[i].Count != second[i].Count {
			t.Errorf("Results() not stable at index %d: %+v vs %+v", i, first[i], second[i])
		}
	}
}
