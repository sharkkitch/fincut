package template

import (
	"testing"
)

func TestNew_EmptyTemplate(t *testing.T) {
	_, err := New(Options{})
	if err == nil {
		t.Fatal("expected error for empty template")
	}
}

func TestNew_InvalidTemplate(t *testing.T) {
	_, err := New(Options{Template: "{{ .Unclosed"})
	if err == nil {
		t.Fatal("expected parse error")
	}
}

func TestNew_InvalidPattern(t *testing.T) {
	_, err := New(Options{Template: "{{.Line}}", Pattern: "[invalid"})
	if err == nil {
		t.Fatal("expected error for invalid pattern")
	}
}

func TestNew_ValidOptions(t *testing.T) {
	_, err := New(Options{Template: "{{.Line}}"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestTemplater_Apply_IdentityTemplate(t *testing.T) {
	tr, _ := New(Options{Template: "{{.Line}}"})
	lines := []string{"hello", "world"}
	out, err := tr.Apply(lines)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for i, l := range lines {
		if out[i] != l {
			t.Errorf("line %d: got %q, want %q", i, out[i], l)
		}
	}
}

func TestTemplater_Apply_StaticTemplate(t *testing.T) {
	tr, _ := New(Options{Template: "prefix: {{.Line}}"})
	out, err := tr.Apply([]string{"msg"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out[0] != "prefix: msg" {
		t.Errorf("got %q", out[0])
	}
}

func TestTemplater_Apply_WithPattern_Match(t *testing.T) {
	tr, _ := New(Options{
		Template: "level={{.level}} msg={{.msg}}",
		Pattern:  `(?P<level>\w+) (?P<msg>.+)`,
	})
	out, err := tr.Apply([]string{"ERROR something went wrong"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := "level=ERROR msg=something went wrong"
	if out[0] != want {
		t.Errorf("got %q, want %q", out[0], want)
	}
}

func TestTemplater_Apply_WithPattern_NoMatch_PassThrough(t *testing.T) {
	tr, _ := New(Options{
		Template: "level={{.level}}",
		Pattern:  `^NOMATCH`,
	})
	line := "INFO something"
	out, err := tr.Apply([]string{line})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out[0] != line {
		t.Errorf("expected pass-through, got %q", out[0])
	}
}

func TestTemplater_Apply_EmptyInput(t *testing.T) {
	tr, _ := New(Options{Template: "{{.Line}}"})
	out, err := tr.Apply([]string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 0 {
		t.Errorf("expected empty output, got %d lines", len(out))
	}
}
