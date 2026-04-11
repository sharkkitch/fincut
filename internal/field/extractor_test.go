package field

import (
	"testing"
)

func TestNew_NeitherPatternNorDelimiter(t *testing.T) {
	_, err := New(Options{Fields: []string{"a"}})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestNew_EmptyFields(t *testing.T) {
	_, err := New(Options{Delimiter: " ", Fields: []string{}})
	if err == nil {
		t.Fatal("expected error for empty Fields")
	}
}

func TestNew_InvalidPattern(t *testing.T) {
	_, err := New(Options{Pattern: "([unclosed", Fields: []string{"x"}})
	if err == nil {
		t.Fatal("expected compile error")
	}
}

func TestNew_PatternMissingCaptureGroup(t *testing.T) {
	_, err := New(Options{Pattern: `(?P<level>\w+)`, Fields: []string{"level", "missing"}})
	if err == nil {
		t.Fatal("expected error for missing capture group")
	}
}

func TestNew_ValidDelimiter(t *testing.T) {
	_, err := New(Options{Delimiter: "|", Fields: []string{"ts", "level", "msg"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestExtract_Delimiter_AllFields(t *testing.T) {
	ex, _ := New(Options{Delimiter: " ", Fields: []string{"ts", "level", "msg"}})
	got := ex.Extract("2024-01-01 INFO hello world")
	if got["ts"] != "2024-01-01" {
		t.Errorf("ts: want 2024-01-01, got %q", got["ts"])
	}
	if got["level"] != "INFO" {
		t.Errorf("level: want INFO, got %q", got["level"])
	}
	if got["msg"] != "hello" {
		t.Errorf("msg: want hello, got %q", got["msg"])
	}
}

func TestExtract_Delimiter_FewerTokensThanFields(t *testing.T) {
	ex, _ := New(Options{Delimiter: "|", Fields: []string{"a", "b", "c"}})
	got := ex.Extract("x|y")
	if got["a"] != "x" || got["b"] != "y" {
		t.Errorf("unexpected values: %v", got)
	}
	if got["c"] != "" {
		t.Errorf("c should be empty, got %q", got["c"])
	}
}

func TestExtract_Pattern_Match(t *testing.T) {
	ex, err := New(Options{
		Pattern: `(?P<level>\w+)\s+(?P<msg>.+)`,
		Fields:  []string{"level", "msg"},
	})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	got := ex.Extract("ERROR something went wrong")
	if got["level"] != "ERROR" {
		t.Errorf("level: want ERROR, got %q", got["level"])
	}
	if got["msg"] != "something went wrong" {
		t.Errorf("msg: want 'something went wrong', got %q", got["msg"])
	}
}

func TestExtract_Pattern_NoMatch(t *testing.T) {
	ex, _ := New(Options{
		Pattern: `(?P<level>\d+)`,
		Fields:  []string{"level"},
	})
	got := ex.Extract("no digits here at all")
	if got["level"] != "" {
		t.Errorf("expected empty string for no match, got %q", got["level"])
	}
}
