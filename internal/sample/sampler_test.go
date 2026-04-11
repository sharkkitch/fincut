package sample

import (
	"testing"
)

func TestNew_MutuallyExclusive(t *testing.T) {
	_, err := New(Options{Rate: 2, Probability: 0.5})
	if err == nil {
		t.Fatal("expected error for both Rate and Probability set")
	}
}

func TestNew_NeitherSet(t *testing.T) {
	_, err := New(Options{})
	if err == nil {
		t.Fatal("expected error when neither Rate nor Probability is set")
	}
}

func TestNew_NegativeRate(t *testing.T) {
	_, err := New(Options{Rate: -1})
	if err == nil {
		t.Fatal("expected error for negative Rate")
	}
}

func TestNew_InvalidProbability(t *testing.T) {
	for _, p := range []float64{-0.1, 1.1} {
		_, err := New(Options{Probability: p})
		if err == nil {
			t.Fatalf("expected error for probability %f", p)
		}
	}
}

func TestNew_ValidRate(t *testing.T) {
	_, err := New(Options{Rate: 3})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNew_ValidProbability(t *testing.T) {
	_, err := New(Options{Probability: 0.5, Seed: 42})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestSampler_Apply_EmptyInput(t *testing.T) {
	s, _ := New(Options{Rate: 2})
	out := s.Apply([]string{})
	if len(out) != 0 {
		t.Fatalf("expected empty output, got %d lines", len(out))
	}
}

func TestSampler_Apply_Rate(t *testing.T) {
	s, _ := New(Options{Rate: 3})
	input := []string{"a", "b", "c", "d", "e", "f", "g"}
	out := s.Apply(input)
	// indices 0, 3, 6 → "a", "d", "g"
	expected := []string{"a", "d", "g"}
	if len(out) != len(expected) {
		t.Fatalf("expected %d lines, got %d", len(expected), len(out))
	}
	for i, v := range expected {
		if out[i] != v {
			t.Errorf("line %d: expected %q, got %q", i, v, out[i])
		}
	}
}

func TestSampler_Apply_Rate1_KeepsAll(t *testing.T) {
	s, _ := New(Options{Rate: 1})
	input := []string{"x", "y", "z"}
	out := s.Apply(input)
	if len(out) != len(input) {
		t.Fatalf("rate=1 should keep all lines, got %d", len(out))
	}
}

func TestSampler_Apply_Probability_Deterministic(t *testing.T) {
	s, _ := New(Options{Probability: 1.0, Seed: 1})
	input := make([]string, 20)
	for i := range input {
		input[i] = "line"
	}
	out := s.Apply(input)
	if len(out) != 20 {
		t.Fatalf("probability=1.0 should keep all lines, got %d", len(out))
	}
}
