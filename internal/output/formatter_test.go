package output

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func TestNewFormatter_InvalidFormat(t *testing.T) {
	_, err := NewFormatter("xml", &bytes.Buffer{})
	if err == nil {
		t.Fatal("expected error for unsupported format, got nil")
	}
}

func TestNewFormatter_ValidFormats(t *testing.T) {
	formats := []Format{FormatPlain, FormatJSON, FormatColor}
	for _, f := range formats {
		_, err := NewFormatter(f, &bytes.Buffer{})
		if err != nil {
			t.Errorf("unexpected error for format %q: %v", f, err)
		}
	}
}

func TestFormatter_Plain(t *testing.T) {
	var buf bytes.Buffer
	fmt, _ := NewFormatter(FormatPlain, &buf)
	if err := fmt.WriteLine("hello world"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := strings.TrimSpace(buf.String()); got != "hello world" {
		t.Errorf("expected %q, got %q", "hello world", got)
	}
}

func TestFormatter_JSON(t *testing.T) {
	var buf bytes.Buffer
	fmt, _ := NewFormatter(FormatJSON, &buf)
	if err := fmt.WriteLine("some log line"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var result map[string]string
	if err := json.Unmarshal([]byte(strings.TrimSpace(buf.String())), &result); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	if result["line"] != "some log line" {
		t.Errorf("expected line %q, got %q", "some log line", result["line"])
	}
}

func TestFormatter_Color_ContainsANSI(t *testing.T) {
	cases := []struct {
		input    string
		contains string
	}{
		{"ERROR: something failed", "\033[31m"},
		{"WARN: low disk space", "\033[33m"},
		{"INFO: server started", "\033[32m"},
		{"DEBUG: verbose output", "\033[36m"},
		{"plain line", "plain line"},
	}
	for _, tc := range cases {
		var buf bytes.Buffer
		f, _ := NewFormatter(FormatColor, &buf)
		_ = f.WriteLine(tc.input)
		if !strings.Contains(buf.String(), tc.contains) {
			t.Errorf("input %q: expected output to contain %q, got %q", tc.input, tc.contains, buf.String())
		}
	}
}
