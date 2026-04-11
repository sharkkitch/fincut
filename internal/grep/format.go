package grep

import "fmt"

// FormatMatch renders a Match to a human-readable string including context.
func FormatMatch(m Match) string {
	var out string
	for i, l := range m.Before {
		lineNo := m.LineNumber - len(m.Before) + i
		out += fmt.Sprintf("%d- %s\n", lineNo, l)
	}
	out += fmt.Sprintf("%d: %s\n", m.LineNumber, m.Line)
	for i, l := range m.After {
		out += fmt.Sprintf("%d- %s\n", m.LineNumber+1+i, l)
	}
	return out
}

// FormatSummary returns a short summary line for a slice of matches.
func FormatSummary(matches []Match, totalLines int) string {
	return fmt.Sprintf("matched %d/%d lines", len(matches), totalLines)
}
