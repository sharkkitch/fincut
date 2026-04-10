package diff

import (
	"fmt"
	"strings"
)

// LineChange represents a single line-level change between two log snapshots.
type LineChange struct {
	Type    ChangeType
	LineNum int
	Content string
}

// ChangeType indicates whether a line was added, removed, or unchanged.
type ChangeType int

const (
	Unchanged ChangeType = iota
	Added
	Removed
)

func (c ChangeType) String() string {
	switch c {
	case Added:
		return "added"
	case Removed:
		return "removed"
	default:
		return "unchanged"
	}
}

// Diff computes a simple line-level diff between two slices of log lines.
// It returns a slice of LineChange entries describing additions and removals.
func Diff(before, after []string) []LineChange {
	beforeSet := make(map[string]int, len(before))
	for i, line := range before {
		beforeSet[line] = i + 1
	}

	afterSet := make(map[string]int, len(after))
	for i, line := range after {
		afterSet[line] = i + 1
	}

	var changes []LineChange

	for i, line := range before {
		if _, found := afterSet[line]; !found {
			changes = append(changes, LineChange{Type: Removed, LineNum: i + 1, Content: line})
		}
	}

	for i, line := range after {
		if _, found := beforeSet[line]; !found {
			changes = append(changes, LineChange{Type: Added, LineNum: i + 1, Content: line})
		}
	}

	return changes
}

// FormatChange returns a human-readable string for a single LineChange.
func FormatChange(c LineChange) string {
	prefix := " "
	switch c.Type {
	case Added:
		prefix = "+"
	case Removed:
		prefix = "-"
	}
	return fmt.Sprintf("%s [line %d] %s", prefix, c.LineNum, strings.TrimRight(c.Content, "\n"))
}
