// Package filter provides regex-based pipeline filtering for log lines.
//
// A Pipeline is an ordered sequence of filter stages, each consisting of
// a compiled regular expression and an invert flag. A log line passes the
// pipeline only when it satisfies every stage in order: a non-inverted stage
// requires the pattern to match, while an inverted stage requires it not to.
//
// Pipelines are constructed via NewPipeline, which accepts a slice of
// pattern strings and a corresponding slice of invert booleans. The two
// slices must have the same length; otherwise NewPipeline returns an error.
//
// An empty pipeline (no stages) matches every line.
//
// Example usage:
//
//	// Keep lines containing "ERROR" but not "timeout".
//	p, err := filter.NewPipeline([]string{`ERROR`, `timeout`}, []bool{false, true})
//	if err != nil {
//		log.Fatal(err)
//	}
//	matched := p.Match("ERROR: connection refused") // true
//	skipped := p.Match("ERROR: timeout")            // false
package filter
