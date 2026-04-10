package stats

import (
	"strings"
	"testing"
	"time"
)

func TestNewCollector(t *testing.T) {
	c := NewCollector(3)
	if c.FilterStages != 3 {
		t.Errorf("expected FilterStages=3, got %d", c.FilterStages)
	}
	if c.LinesRead != 0 || c.LinesMatched != 0 || c.LinesDropped != 0 {
		t.Error("expected zero counters on creation")
	}
	if c.StartTime.IsZero() {
		t.Error("expected StartTime to be set")
	}
}

func TestCollector_Record_Matched(t *testing.T) {
	c := NewCollector(1)
	c.Record("hello world", true)
	if c.LinesRead != 1 {
		t.Errorf("expected LinesRead=1, got %d", c.LinesRead)
	}
	if c.LinesMatched != 1 {
		t.Errorf("expected LinesMatched=1, got %d", c.LinesMatched)
	}
	if c.LinesDropped != 0 {
		t.Errorf("expected LinesDropped=0, got %d", c.LinesDropped)
	}
	if c.BytesRead != int64(len("hello world")+1) {
		t.Errorf("unexpected BytesRead: %d", c.BytesRead)
	}
}

func TestCollector_Record_Dropped(t *testing.T) {
	c := NewCollector(1)
	c.Record("ignored line", false)
	if c.LinesDropped != 1 {
		t.Errorf("expected LinesDropped=1, got %d", c.LinesDropped)
	}
	if c.LinesMatched != 0 {
		t.Errorf("expected LinesMatched=0, got %d", c.LinesMatched)
	}
}

func TestCollector_MatchRate(t *testing.T) {
	c := NewCollector(0)
	if c.MatchRate() != 0 {
		t.Error("expected MatchRate=0 when no lines read")
	}
	c.Record("a", true)
	c.Record("b", false)
	rate := c.MatchRate()
	if rate < 0.49 || rate > 0.51 {
		t.Errorf("expected MatchRate~0.5, got %f", rate)
	}
}

func TestCollector_Elapsed(t *testing.T) {
	c := NewCollector(0)
	time.Sleep(5 * time.Millisecond)
	if c.Elapsed() < 5*time.Millisecond {
		t.Error("expected elapsed >= 5ms")
	}
}

func TestCollector_Summary(t *testing.T) {
	c := NewCollector(2)
	c.Record("line one", true)
	c.Record("line two", false)
	summary := c.Summary()
	for _, want := range []string{"lines read: 2", "lines matched: 1", "lines dropped: 1", "filter stages: 2", "match rate:", "elapsed:"} {
		if !strings.Contains(summary, want) {
			t.Errorf("summary missing %q\ngot:\n%s", want, summary)
		}
	}
}
