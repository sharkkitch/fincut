package summarize

import (
	"strings"
	"testing"
)

func TestNew_NegativeTopN(t *testing.T) {
	_, err := New(Options{TopN: -1})
	if err == nil {
		t.Fatal("expected error for negative TopN")
	}
}

func TestNew_NegativeMinCount(t *testing.T) {
	_, err := New(Options{MinCount: -3})
	if err == nil {
		t.Fatal("expected error for negative MinCount")
	}
}

func TestNew_ValidOptions(t *testing.T) {
	s, err := New(Options{TopN: 5, MinCount: 2, CaseSensitive: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s == nil {
		t.Fatal("expected non-nil Summarizer")
	}
}

func TestSummarizer_Total(t *testing.T) {
	s, _ := New(Options{})
	for i := 0; i < 7; i++ {
		s.Add("line")
	}
	if s.Total() != 7 {
		t.Fatalf("want 7, got %d", s.Total())
	}
}

func TestSummarizer_Results_Sorted(t *testing.T) {
	s, _ := New(Options{})
	for _, l := range []string{"a", "b", "a", "c", "a", "b"} {
		s.Add(l)
	}
	res := s.Results()
	if len(res) != 3 {
		t.Fatalf("want 3 results, got %d", len(res))
	}
	if res[0].Line != "a" || res[0].Count != 3 {
		t.Fatalf("want a=3 first, got %s=%d", res[0].Line, res[0].Count)
	}
	if res[1].Line != "b" || res[1].Count != 2 {
		t.Fatalf("want b=2 second, got %s=%d", res[1].Line, res[1].Count)
	}
}

func TestSummarizer_TopN(t *testing.T) {
	s, _ := New(Options{TopN: 2})
	for _, l := range []string{"x", "y", "z", "x", "x", "y"} {
		s.Add(l)
	}
	res := s.Results()
	if len(res) != 2 {
		t.Fatalf("want 2 results, got %d", len(res))
	}
}

func TestSummarizer_MinCount(t *testing.T) {
	s, _ := New(Options{MinCount: 3})
	for _, l := range []string{"a", "a", "a", "b", "b", "c"} {
		s.Add(l)
	}
	res := s.Results()
	if len(res) != 1 || res[0].Line != "a" {
		t.Fatalf("expected only 'a', got %v", res)
	}
}

func TestSummarizer_CaseInsensitive(t *testing.T) {
	s, _ := New(Options{CaseSensitive: false})
	for _, l := range []string{"ERROR", "error", "Error"} {
		s.Add(l)
	}
	res := s.Results()
	if len(res) != 1 {
		t.Fatalf("want 1 distinct entry, got %d", len(res))
	}
	if res[0].Count != 3 {
		t.Fatalf("want count 3, got %d", res[0].Count)
	}
}

func TestFormatSummary_ContainsLine(t *testing.T) {
	res := []Result{{Line: "foo", Count: 42}}
	out := FormatSummary(res, 100)
	if !strings.Contains(out, "foo") {
		t.Fatalf("expected 'foo' in output: %s", out)
	}
	if !strings.Contains(out, "42") {
		t.Fatalf("expected '42' in output: %s", out)
	}
}
