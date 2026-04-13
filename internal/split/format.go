package split

import "fmt"

// FormatSummary returns a human-readable summary of the groups produced by the
// Splitter.
func FormatSummary(groups []Group) string {
	total := 0
	for _, g := range groups {
		total += len(g.Lines)
	}
	return fmt.Sprintf("%d group(s), %d line(s) total", len(groups), total)
}

// FormatGroup renders a single group with its label header.
func FormatGroup(g Group) string {
	header := fmt.Sprintf("=== %s (%d lines) ===", g.Label, len(g.Lines))
	if len(g.Lines) == 0 {
		return header
	}
	out := header
	for _, l := range g.Lines {
		out += "\n" + l
	}
	return out
}
