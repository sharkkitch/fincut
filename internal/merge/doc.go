// Package merge provides utilities for combining log lines from multiple
// sources into a single ordered stream.
//
// # Overview
//
// A Merger accepts a map of source names to line slices and writes the
// combined output to any io.Writer. Two ordering modes are supported:
//
//   - Default (insertion order per source, sources iterated in map order)
//   - SortByTime: lines are re-ordered by a parsed timestamp prefix using
//     the Go time layout specified in Options.TimestampLayout.
//
// # Source Labelling
//
// When Options.LabelSources is true every output line is prefixed with
// the source filename enclosed in square brackets, e.g.:
//
//	[app.log] 2024-06-01T12:00:01 server started
//
// This is useful when merging logs from heterogeneous services and you
// need to trace each line back to its origin.
//
// # Error Handling
//
// If SortByTime is enabled and a line does not begin with a recognisable
// timestamp (according to Options.TimestampLayout), the line is placed
// before all timestamped lines in its source group. A non-nil error is
// returned after writing so the caller can decide whether to treat
// unparseable lines as a fatal condition.
package merge
