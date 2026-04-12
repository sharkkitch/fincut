package parse

import (
	"strings"
	"testing"
)

func TestNew_UnknownFormat(t *testing.T) {
	_, err := New(Options{Format: "csv"})
	if err == nil {
		t.Fatal("expected error for unknown format")
	}
}

func TestNew_RegexMissingPattern(t *testing.T) {
	_, err := New(Options{Format: FormatRegex, Pattern: ""})
	if err == nil {
		t.Fatal("expected error for missing pattern")
	}
}

func TestNew_RegexInvalidPattern(t *testing.T) {
	_, err := New(Options{Format: FormatRegex, Pattern: "[invalid"})
	if err == nil {
		t.Fatal("expected error for invalid regex")
	}
}

func TestNew_DelimMissingDelimiter(t *testing.T) {
	_, err := New(Options{Format: FormatDelim, Fields: []string{"a"}})
	if err == nil {
		t.Fatal("expected error for missing delimiter")
	}
}

func TestNew_DelimMissingFields(t *testing.T) {
	_, err := New(Options{Format: FormatDelim, Delimiter: "|"})
	if err == nil {
		t.Fatal("expected error for missing fields")
	}
}

func TestParser_JSON_Valid(t *testing.T) {
	p, err := New(Options{Format: FormatJSON})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	fields, err := p.Parse(`{"level":"info","msg":"started","pid":"42"}`)
	if err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}
	if fields["level"] != "info" {
		t.Errorf("expected level=info, got %q", fields["level"])
	}
	if fields["msg"] != "started" {
		t.Errorf("expected msg=started, got %q", fields["msg"])
	}
}

func TestParser_JSON_Invalid(t *testing.T) {
	p, _ := New(Options{Format: FormatJSON})
	_, err := p.Parse("not json")
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestParser_Regex_Match(t *testing.T) {
	p, err := New(Options{
		Format:  FormatRegex,
		Pattern: `(?P<level>\w+)\s+(?P<msg>.+)`,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	fields, err := p.Parse("ERROR something went wrong")
	if err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}
	if fields["level"] != "ERROR" {
		t.Errorf("expected level=ERROR, got %q", fields["level"])
	}
	if !strings.HasPrefix(fields["msg"], "something") {
		t.Errorf("unexpected msg: %q", fields["msg"])
	}
}

func TestParser_Regex_NoMatch(t *testing.T) {
	p, _ := New(Options{Format: FormatRegex, Pattern: `(?P<ts>\d{4}-\d{2}-\d{2})`})
	_, err := p.Parse("no date here")
	if err == nil {
		t.Fatal("expected error for non-matching line")
	}
}

func TestParser_Delim_Valid(t *testing.T) {
	p, err := New(Options{
		Format:    FormatDelim,
		Delimiter: "|",
		Fields:    []string{"ts", "level", "msg"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	fields, err := p.Parse("2024-01-01|INFO|hello world")
	if err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}
	if fields["level"] != "INFO" {
		t.Errorf("expected level=INFO, got %q", fields["level"])
	}
	if fields["msg"] != "hello world" {
		t.Errorf("expected msg='hello world', got %q", fields["msg"])
	}
}

func TestParser_Delim_ShortLine(t *testing.T) {
	p, _ := New(Options{
		Format:    FormatDelim,
		Delimiter: ",",
		Fields:    []string{"a", "b", "c"},
	})
	fields, err := p.Parse("x,y")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if fields["c"] != "" {
		t.Errorf("expected empty string for missing field c, got %q", fields["c"])
	}
}
