// Package window provides sliding and fixed-size window operations over
// slices of log lines.
//
// A Windower partitions an input slice into consecutive (or overlapping)
// sub-slices of a fixed size. This is useful for batch processing, rolling
// analysis, or chunked diffing of log streams.
//
// Basic usage:
//
//	w, err := window.New(window.Options{Size: 10, Step: 5})
//	if err != nil { ... }
//	windows := w.Apply(lines)
//
// Use Flatten to reconstruct the original line order from a set of windows,
// automatically removing duplicates introduced by overlap.
package window
