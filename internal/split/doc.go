// Package split provides a line-stream splitter that partitions input into
// named groups whenever a configurable regex delimiter is matched.
//
// Basic usage:
//
//	s, err := split.New(split.Options{
//		Pattern:       `^---`,
//		KeepDelimiter: false,
//		Label:         "section",
//	})
//	if err != nil { ... }
//	groups := s.Apply(lines)
//
// Each group is a [Group] value containing a generated label and the lines
// that belong to it.  The delimiter line itself is discarded unless
// KeepDelimiter is true.
package split
