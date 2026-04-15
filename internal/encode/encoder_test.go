package encode

import (
	"strings"
	"testing"
)

func TestNew_EmptyEncoding(t *testing.T) {
	_, err := New(Options{})
	if err == nil {
		t.Fatal("expected error for empty encoding, got nil")
	}
}

func TestNew_UnknownEncoding(t *testing.T) {
	_, err := New(Options{Encoding: "rot13"})
	if err == nil {
		t.Fatal("expected error for unknown encoding, got nil")
	}
	if !strings.Contains(err.Error(), "unknown encoding") {
		t.Errorf("expected 'unknown encoding' in error, got: %v", err)
	}
}

func TestNew_ValidOptions(t *testing.T) {
	for _, enc := range []Encoding{Base64, Base64URL, Hex} {
		_, err := New(Options{Encoding: enc})
		if err != nil {
			t.Errorf("unexpected error for encoding %q: %v", enc, err)
		}
	}
}

func TestEncoder_Apply_Base64RoundTrip(t *testing.T) {
	input := []string{"hello world", "log line: INFO started"}

	enc, _ := New(Options{Encoding: Base64})
	encoded, err := enc.Apply(input)
	if err != nil {
		t.Fatalf("encode error: %v", err)
	}

	dec, _ := New(Options{Encoding: Base64, Decode: true})
	decoded, err := dec.Apply(encoded)
	if err != nil {
		t.Fatalf("decode error: %v", err)
	}

	for i, line := range decoded {
		if line != input[i] {
			t.Errorf("line %d: got %q, want %q", i, line, input[i])
		}
	}
}

func TestEncoder_Apply_HexRoundTrip(t *testing.T) {
	input := []string{"deadbeef", "structured log"}

	enc, _ := New(Options{Encoding: Hex})
	encoded, err := enc.Apply(input)
	if err != nil {
		t.Fatalf("encode error: %v", err)
	}

	dec, _ := New(Options{Encoding: Hex, Decode: true})
	decoded, err := dec.Apply(encoded)
	if err != nil {
		t.Fatalf("decode error: %v", err)
	}

	for i, line := range decoded {
		if line != input[i] {
			t.Errorf("line %d: got %q, want %q", i, line, input[i])
		}
	}
}

func TestEncoder_Apply_Base64URLRoundTrip(t *testing.T) {
	input := []string{"url-safe content", "another line"}

	enc, _ := New(Options{Encoding: Base64URL})
	encoded, err := enc.Apply(input)
	if err != nil {
		t.Fatalf("encode error: %v", err)
	}

	dec, _ := New(Options{Encoding: Base64URL, Decode: true})
	decoded, err := dec.Apply(encoded)
	if err != nil {
		t.Fatalf("decode error: %v", err)
	}

	for i, line := range decoded {
		if line != input[i] {
			t.Errorf("line %d: got %q, want %q", i, line, input[i])
		}
	}
}

func TestEncoder_Decode_InvalidBase64(t *testing.T) {
	dec, _ := New(Options{Encoding: Base64, Decode: true})
	_, err := dec.Apply([]string{"!!!not-valid-base64!!!"})
	if err == nil {
		t.Fatal("expected error for invalid base64 input, got nil")
	}
}

func TestEncoder_Apply_EmptyInput(t *testing.T) {
	enc, _ := New(Options{Encoding: Hex})
	out, err := enc.Apply([]string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 0 {
		t.Errorf("expected empty output, got %d lines", len(out))
	}
}
