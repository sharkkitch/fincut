// Package burst provides log burst detection over sliding time windows.
//
// A Burster scans a slice of log lines, extracts timestamps using a
// configurable regex pattern and time layout, and groups consecutive lines
// that exceed a caller-defined rate threshold (lines per second) within a
// given window duration.
//
// Example usage:
//
//	b, err := burst.New(burst.Options{
//		TimestampPattern: `\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}`,
//		TimestampLayout:  "2006-01-02T15:04:05",
//		Window:           5 * time.Second,
//		Threshold:        10.0,
//	})
//	if err != nil { ... }
//	bursts, err := b.Apply(lines)
package burst
