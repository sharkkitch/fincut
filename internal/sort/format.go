package sort

import "fmt"

// FormatSummary returns a human-readable summary of a sort operation.
func FormatSummary(inputCount, outputCount int, reverse, unique bool) string {
	direction := "ascending"
	if reverse {
		direction = "descending"
	}
	if unique {
		return fmt.Sprintf("sorted %d lines %s → %d unique lines", inputCount, direction, outputCount)
	}
	return fmt.Sprintf("sorted %d lines %s", inputCount, direction)
}
