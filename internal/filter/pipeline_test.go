package filter

import (
	"testing"
)

func TestNewPipeline_InvalidPattern(t *testing.T) {
	_, err := NewPipeline([]string{"["})
	if err == nil {
		t.Fatal("expected error for invalid regex, got nil")
	}
}

func TestPipeline_Match_SingleInclude(t *testing.T) {
	p, err := NewPipeline([]string{"ERROR"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !p.Match("2024-01-01 ERROR something broke") {
		t.Error("expected match for line containing ERROR")
	}
	if p.Match("2024-01-01 INFO all good") {
		t.Error("expected no match for line without ERROR")
	}
}

func TestPipeline_Match_InvertedPattern(t *testing.T) {
	p, err := NewPipeline([]string{"!DEBUG"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.Match("2024-01-01 DEBUG verbose output") {
		t.Error("expected no match for DEBUG line with inverted filter")
	}
	if !p.Match("2024-01-01 INFO startup complete") {
		t.Error("expected match for non-DEBUG line with inverted filter")
	}
}

func TestPipeline_Match_MultiStage(t *testing.T) {
	p, err := NewPipeline([]string{"ERROR", "!timeout"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !p.Match("ERROR connection refused") {
		t.Error("expected match: ERROR without timeout")
	}
	if p.Match("ERROR timeout exceeded") {
		t.Error("expected no match: ERROR with timeout")
	}
	if p.Match("INFO all good") {
		t.Error("expected no match: no ERROR")
	}
}

func TestPipeline_Apply(t *testing.T) {
	lines := []string{
		"INFO server started",
		"ERROR disk full",
		"ERROR timeout",
		"DEBUG heartbeat",
	}
	p, err := NewPipeline([]string{"ERROR", "!timeout"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	result := p.Apply(lines)
	if len(result) != 1 || result[0] != "ERROR disk full" {
		t.Errorf("unexpected result: %v", result)
	}
}

func TestPipeline_EmptyPipeline(t *testing.T) {
	p, err := NewPipeline(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	lines := []string{"anything", "goes"}
	result := p.Apply(lines)
	if len(result) != 2 {
		t.Errorf("empty pipeline should pass all lines, got %v", result)
	}
}
