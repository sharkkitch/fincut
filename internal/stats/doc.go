// Package stats provides a lightweight collector for tracking fincut
// processing metrics, including line counts, byte throughput, filter
// match rates, and elapsed time.
//
// Usage:
//
//	col := stats.NewCollector(len(pipeline.Stages))
//	for _, line := range lines {
//		matched := pipeline.Match(line)
//		col.Record(line, matched)
//	}
//	fmt.Println(col.Summary())
package stats
