package pivot

import (
	"strings"
	"testing"
)

func TestNew_InvalidKeyField(t *testing.T) {
	_, err := New(Options{KeyField: -1, Delimiter: ","})
	if err == nil {
		t.Fatal("expected error for negative key field")
	}
}

func TestNew_EmptyDelimiter(t *testing.T) {
	_, err := New(Options{KeyField: 0, Delimiter: ""})
	if err == nil {
		t.Fatal("expected error for empty delimiter")
	}
}

func TestNew_UnknownAggregator(t *testing.T) {
	_, err := New(Options{KeyField: 0, Delimiter: ",", Aggregator: "median"})
	if err == nil {
		t.Fatal("expected error for unknown aggregator")
	}
}

func TestNew_ValidOptions(t *testing.T) {
	p, err := New(Options{KeyField: 0, ValueField: 1, Delimiter: ",", Aggregator: "count"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p == nil {
		t.Fatal("expected non-nil Pivotter")
	}
}

func TestPivotter_Apply_Count(t *testing.T) {
	p, _ := New(Options{KeyField: 0, ValueField: 1, Delimiter: ",", Aggregator: "count"})
	lines := []string{"a,1", "b,2", "a,3", "b,4", "c,5"}
	out := p.Apply(lines)
	if len(out) != 3 {
		t.Fatalf("expected 3 groups, got %d", len(out))
	}
	if out[0] != "a,2" {
		t.Errorf("expected a,2 got %s", out[0])
	}
}

func TestPivotter_Apply_Sum(t *testing.T) {
	p, _ := New(Options{KeyField: 0, ValueField: 1, Delimiter: ":", Aggregator: "sum"})
	lines := []string{"x:10", "x:5", "y:3"}
	out := p.Apply(lines)
	if len(out) != 2 {
		t.Fatalf("expected 2 groups, got %d", len(out))
	}
	if out[0] != "x:15" {
		t.Errorf("expected x:15 got %s", out[0])
	}
}

func TestPivotter_Apply_Values(t *testing.T) {
	p, _ := New(Options{KeyField: 0, ValueField: 1, Delimiter: ",", Aggregator: "values"})
	lines := []string{"k,alpha", "k,beta"}
	out := p.Apply(lines)
	if len(out) != 1 {
		t.Fatalf("expected 1 group, got %d", len(out))
	}
	if !strings.Contains(out[0], "alpha") || !strings.Contains(out[0], "beta") {
		t.Errorf("unexpected output: %s", out[0])
	}
}

func TestPivotter_Apply_SkipsOutOfBoundsKey(t *testing.T) {
	p, _ := New(Options{KeyField: 3, ValueField: -1, Delimiter: ",", Aggregator: "count"})
	lines := []string{"only,two,fields", "a,b,c,d"}
	out := p.Apply(lines)
	if len(out) != 1 {
		t.Fatalf("expected 1 result, got %d", len(out))
	}
}

func TestPivotter_Apply_EmptyInput(t *testing.T) {
	p, _ := New(Options{KeyField: 0, ValueField: 1, Delimiter: ","})
	out := p.Apply([]string{})
	if len(out) != 0 {
		t.Fatalf("expected empty output, got %d lines", len(out))
	}
}
