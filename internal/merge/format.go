package merge

import "fmt"

// FormatSummary returns a human-readable summary of a completed merge operation.
func FormatSummary(sourceCounts map[string]int, total int) string {
	s := fmt.Sprintf("merge: %d total lines from %d source(s)\n", total, len(sourceCounts))
	for src, count := range sourceCounts {
		s += fmt.Sprintf("  %-30s %d lines\n", src, count)
	}
	return s
}

// CountLines tallies the number of lines contributed by each source.
func CountLines(sources map[string][]string) (map[string]int, int) {
	counts := make(map[string]int, len(sources))
	total := 0
	for src, lines := range sources {
		counts[src] = len(lines)
		total += len(lines)
	}
	return counts, total
}
