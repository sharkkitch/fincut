package segment

import (
	"fmt"
	"strings"
)

// FormatSegment renders a Segment as a human-readable block.
// Each line is prefixed with the segment label and a line index.
func FormatSegment(s Segment) string {
	if len(s.Lines) == 0 {
		return fmt.Sprintf("[%s] (empty)\n", s.Label)
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("[%s] %d lines\n", s.Label, len(s.Lines)))
	for i, line := range s.Lines {
		sb.WriteString(fmt.Sprintf("  %d: %s\n", i+1, line))
	}
	return sb.String()
}

// SummaryTable returns a compact multi-segment overview string suitable
// for plain-text reporting.
func SummaryTable(segs []Segment) string {
	if len(segs) == 0 {
		return "no segments\n"
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%-20s %s\n", "SEGMENT", "LINES"))
	sb.WriteString(strings.Repeat("-", 30) + "\n")
	for _, s := range segs {
		sb.WriteString(fmt.Sprintf("%-20s %d\n", s.Label, len(s.Lines)))
	}
	return sb.String()
}
