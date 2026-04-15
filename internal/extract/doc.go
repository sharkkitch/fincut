// Package extract implements regex capture-group extraction for log lines.
//
// Given a pattern with one or more capture groups, Extractor replaces each
// matching line with the text of the selected group. Lines that do not match
// are either passed through unchanged or dropped, depending on SkipUnmatched.
//
// Example usage:
//
//	e, err := extract.New(extract.Options{
//		Pattern:       `level=(\w+)`,
//		Group:         1,
//		SkipUnmatched: true,
//	})
//	result := e.Apply(lines)
package extract
