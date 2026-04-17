// Package between provides a Betweener that extracts lines located between
// two regex boundary patterns. It supports multiple non-overlapping regions
// within the input and an inclusive mode that retains the boundary lines
// themselves in the output.
//
// Basic usage:
//
//	b, err := between.New(between.Options{
//		StartPattern: `^START`,
//		EndPattern:   `^END`,
//		Inclusive:    false,
//	})
//	if err != nil { ... }
//	out := b.Apply(lines)
package between
