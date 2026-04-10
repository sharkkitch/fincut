package stats_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/yourorg/fincut/internal/stats"
)

func TestNewReporter_NilWriter(t *testing.T) {
	_, err := stats.NewReporter(nil, "plain")
	if err == nil {
		t.Fatal("expected error for nil writer, got nil")
	}
}

func TestNewReporter_InvalidFormat(t *testing.T) {
	var buf bytes.Buffer
	_, err := stats.NewReporter(&buf, "xml")
	if err == nil {
		t.Fatal("expected error for invalid format, got nil")
	}
}

func TestReporter_Plain_Output(t *testing.T) {
	var buf bytes.Buffer
	r, err := stats.NewReporter(&buf, "plain")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	c, _ := stats.NewCollector()
	c.Record(true)
	c.Record(true)
	c.Record(false)
	time.Sleep(1 * time.Millisecond)

	r.Report(c)

	out := buf.String()
	if !strings.Contains(out, "matched") {
		t.Errorf("expected 'matched' in plain output, got: %s", out)
	}
	if !strings.Contains(out, "dropped") {
		t.Errorf("expected 'dropped' in plain output, got: %s", out)
	}
}

func TestReporter_JSON_Output(t *testing.T) {
	var buf bytes.Buffer
	r, err := stats.NewReporter(&buf, "json")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	c, _ := stats.NewCollector()
	c.Record(true)
	c.Record(false)
	c.Record(false)

	r.Report(c)

	var result map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("expected valid JSON output, got error: %v\noutput: %s", err, buf.String())
	}

	if _, ok := result["matched"]; !ok {
		t.Error("expected 'matched' key in JSON output")
	}
	if _, ok := result["dropped"]; !ok {
		t.Error("expected 'dropped' key in JSON output")
	}
	if _, ok := result["match_rate"]; !ok {
		t.Error("expected 'match_rate' key in JSON output")
	}
}
