package reorder_test

import (
	"testing"

	"github.com/user/fincut/internal/reorder"
)

func TestNew_NeitherSet(t *testing.T) {
	_, err := reorder.New(reorder.Options{})
	if err == nil {
		t.Fatal("expected error when neither reverse nor field is set")
	}
}

func TestNew_FieldWithoutDelimiter(t *testing.T) {
	_, err := reorder.New(reorder.Options{Field: 2})
	if err == nil {
		t.Fatal("expected error when field set without delimiter")
	}
}

func TestNew_NegativeField(t *testing.T) {
	_, err := reorder.New(reorder.Options{Field: -1, Delimiter: ","})
	if err == nil {
		t.Fatal("expected error for negative field")
	}
}

func TestNew_ValidReverse(t *testing.T) {
	r, err := reorder.New(reorder.Options{Reverse: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r == nil {
		t.Fatal("expected non-nil Reorderer")
	}
}

func TestReorderer_Apply_Reverse(t *testing.T) {
	r, _ := reorder.New(reorder.Options{Reverse: true})
	input := []string{"a", "b", "c"}
	got := r.Apply(input)
	want := []string{"c", "b", "a"}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("index %d: got %q want %q", i, got[i], want[i])
		}
	}
}

func TestReorderer_Apply_SortByField(t *testing.T) {
	r, _ := reorder.New(reorder.Options{Field: 2, Delimiter: ","})
	input := []string{"z,banana", "a,apple", "m,mango"}
	got := r.Apply(input)
	want := []string{"a,apple", "z,banana", "m,mango"}
	for i, w := range want {
		if got[i] != w {
			t.Errorf("index %d: got %q want %q", i, got[i], w)
		}
	}
}

func TestReorderer_Apply_ReverseByField(t *testing.T) {
	r, _ := reorder.New(reorder.Options{Field: 1, Delimiter: ":", Reverse: true})
	input := []string{"beta:1", "alpha:2", "gamma:3"}
	got := r.Apply(input)
	if got[0] != "gamma:3" {
		t.Errorf("expected gamma first, got %q", got[0])
	}
}

func TestReorderer_Apply_DoesNotMutateInput(t *testing.T) {
	r, _ := reorder.New(reorder.Options{Reverse: true})
	input := []string{"x", "y", "z"}
	orig := make([]string, len(input))
	copy(orig, input)
	r.Apply(input)
	for i := range orig {
		if input[i] != orig[i] {
			t.Errorf("input mutated at index %d", i)
		}
	}
}

func TestFormatSummary(t *testing.T) {
	in := []string{"a", "b"}
	out := []string{"b", "a"}
	s := reorder.FormatSummary(in, out)
	if s == "" {
		t.Error("expected non-empty summary")
	}
}
