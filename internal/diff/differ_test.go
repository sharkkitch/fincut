package diff

import (
	"strings"
	"testing"
)

func TestDiff_NoChanges(t *testing.T) {
	lines := []string{"alpha", "beta", "gamma"}
	changes := Diff(lines, lines)
	if len(changes) != 0 {
		t.Errorf("expected 0 changes, got %d", len(changes))
	}
}

func TestDiff_AllAdded(t *testing.T) {
	after := []string{"new line 1", "new line 2"}
	changes := Diff(nil, after)
	if len(changes) != 2 {
		t.Fatalf("expected 2 changes, got %d", len(changes))
	}
	for _, c := range changes {
		if c.Type != Added {
			t.Errorf("expected Added, got %v", c.Type)
		}
	}
}

func TestDiff_AllRemoved(t *testing.T) {
	before := []string{"old line 1", "old line 2"}
	changes := Diff(before, nil)
	if len(changes) != 2 {
		t.Fatalf("expected 2 changes, got %d", len(changes))
	}
	for _, c := range changes {
		if c.Type != Removed {
			t.Errorf("expected Removed, got %v", c.Type)
		}
	}
}

func TestDiff_MixedChanges(t *testing.T) {
	before := []string{"keep", "remove me"}
	after := []string{"keep", "add me"}
	changes := Diff(before, after)
	if len(changes) != 2 {
		t.Fatalf("expected 2 changes, got %d", len(changes))
	}
	types := map[ChangeType]bool{}
	for _, c := range changes {
		types[c.Type] = true
	}
	if !types[Added] || !types[Removed] {
		t.Errorf("expected both Added and Removed changes")
	}
}

func TestFormatChange_Added(t *testing.T) {
	c := LineChange{Type: Added, LineNum: 3, Content: "new entry"}
	out := FormatChange(c)
	if !strings.HasPrefix(out, "+") {
		t.Errorf("expected '+' prefix, got: %s", out)
	}
	if !strings.Contains(out, "line 3") {
		t.Errorf("expected line number in output, got: %s", out)
	}
	if !strings.Contains(out, "new entry") {
		t.Errorf("expected content in output, got: %s", out)
	}
}

func TestFormatChange_Removed(t *testing.T) {
	c := LineChange{Type: Removed, LineNum: 1, Content: "old entry"}
	out := FormatChange(c)
	if !strings.HasPrefix(out, "-") {
		t.Errorf("expected '-' prefix, got: %s", out)
	}
}

func TestChangeType_String(t *testing.T) {
	if Added.String() != "added" {
		t.Errorf("unexpected string for Added")
	}
	if Removed.String() != "removed" {
		t.Errorf("unexpected string for Removed")
	}
	if Unchanged.String() != "unchanged" {
		t.Errorf("unexpected string for Unchanged")
	}
}
