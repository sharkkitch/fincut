package window

import "fmt"

// FormatSummary returns a human-readable summary of windowing results.
func FormatSummary(totalLines, windowSize, step int) string {
	count := 0
	for start := 0; start < totalLines; start += step {
		count++
		end := start + windowSize
		if end >= totalLines {
			break
		}
	}
	if totalLines == 0 {
		count = 0
	}
	return fmt.Sprintf("windows=%d size=%d step=%d total_lines=%d", count, windowSize, step, totalLines)
}

// FormatWindow returns a labelled header for a single window.
func FormatWindow(index, start, end int) string {
	return fmt.Sprintf("[window %d | lines %d-%d]", index+1, start+1, end)
}
