package burst

import "fmt"

// FormatSummary returns a human-readable summary of detected bursts.
func FormatSummary(bursts []Burst) string {
	if len(bursts) == 0 {
		return "no bursts detected"
	}
	return fmt.Sprintf("%d burst(s) detected", len(bursts))
}

// FormatBurst returns a single-line description of a Burst.
func FormatBurst(b Burst) string {
	return fmt.Sprintf(
		"burst lines %d-%d (%d lines, %.2f lines/sec)",
		b.Start, b.End, len(b.Lines), b.Rate,
	)
}
