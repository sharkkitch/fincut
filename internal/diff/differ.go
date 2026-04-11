package diff

import (
	"fmt"
	"strings"
)

// ChangeType represents the kind of change in a diff.
type ChangeType int

const (
	// Unchanged indicates the line exists in both sequences.
	Unchanged ChangeType = iota
	// Added indicates the line was added in the new sequence.
	Added
	// Removed indicates the line was removed from the old sequence.
	Removed
)

// Change represents a single line-level diff entry.
type Change struct {
	Type ChangeType
	Line string
}

// Diff computes a line-level diff between two slices of strings using a
// simple LCS-based algorithm. It returns an ordered slice of Change entries.
func Diff(before, after []string) []Change {
	m := len(before)
	n := len(after)

	// Build LCS table.
	table := make([][]int, m+1)
	for i := range table {
		table[i] = make([]int, n+1)
	}
	for i := 1; i <= m; i++ {
		for j := 1; j <= n; j++ {
			if before[i-1] == after[j-1] {
				table[i][j] = table[i-1][j-1] + 1
			} else if table[i-1][j] >= table[i][j-1] {
				table[i][j] = table[i-1][j]
			} else {
				table[i][j] = table[i][j-1]
			}
		}
	}

	// Backtrack to build change list.
	var changes []Change
	i, j := m, n
	for i > 0 || j > 0 {
		switch {
		case i > 0 && j > 0 && before[i-1] == after[j-1]:
			changes = append([]Change{{Type: Unchanged, Line: before[i-1]}}, changes...)
			i--
			j--
		case j > 0 && (i == 0 || table[i][j-1] >= table[i-1][j]):
			changes = append([]Change{{Type: Added, Line: after[j-1]}}, changes...)
			j--
		default:
			changes = append([]Change{{Type: Removed, Line: before[i-1]}}, changes...)
			i--
		}
	}
	return changes
}

// FormatChange renders a single Change as a diff-style string.
// Added lines are prefixed with "+", removed with "-", unchanged with " ".
func FormatChange(c Change) string {
	switch c.Type {
	case Added:
		return fmt.Sprintf("+ %s", strings.TrimRight(c.Line, "\n"))
	case Removed:
		return fmt.Sprintf("- %s", strings.TrimRight(c.Line, "\n"))
	default:
		return fmt.Sprintf("  %s", strings.TrimRight(c.Line, "\n"))
	}
}
