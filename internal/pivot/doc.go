// Package pivot provides a field-based pivot operation for structured log lines.
//
// A Pivotter groups input lines by a configurable key field and aggregates
// values from a second field using one of the built-in aggregators:
//
//   - count  – number of lines per key (default)
//   - sum    – numeric sum of the value field per key
//   - values – comma-joined list of raw values per key
//
// Fields are split by a caller-supplied delimiter string, making pivot
// compatible with CSV, TSV, or any single-character-separated log format.
//
// Example:
//
//	p, err := pivot.New(pivot.Options{
//		KeyField:   0,
//		ValueField: 1,
//		Delimiter:  ",",
//		Aggregator: "sum",
//	})
//	result := p.Apply(lines)
package pivot
