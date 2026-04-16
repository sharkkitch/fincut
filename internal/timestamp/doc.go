// Package timestamp provides a Stamper that prepends or appends a formatted
// timestamp to every line of input.
//
// Usage:
//
//	s, err := timestamp.New(timestamp.Options{
//		Format:  time.RFC3339,
//		Prepend: true,
//	})
//	if err != nil { ... }
//	out := s.Apply(lines)
//
// The Format field accepts any layout string understood by time.Format.
// Sep controls the separator inserted between the timestamp and the original
// line content (default: single space).
package timestamp
