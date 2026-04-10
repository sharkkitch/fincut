// Package filter provides regex-based pipeline filtering for log lines.
//
// A Pipeline is an ordered sequence of filter stages, each consisting of
// a compiled regular expression and an invert flag. A log line passes the
// pipeline only when it satisfies every stage in order: a non-inverted stage
// requires the pattern to match, while an inverted stage requires it not to.
//
// Pipelines are constructed via NewPipeline, which accepts a slice of
// pattern strings and a corresponding slice of invert booleans.
//
// Example usage:
//
//	p, err := filter.NewPipeline([]string{`ERROR`, `timeout`}, []bool{false, true})
//	if err != nil {
//		log.Fatal(err)
//	}
//	matched := p.Match("ERROR: connection refused")
package filter
