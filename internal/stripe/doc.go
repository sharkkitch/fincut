// Package stripe provides line sampling by periodic selection.
//
// A Striper selects one line out of every N lines, starting at a
// configurable zero-based offset within each group. This is useful
// for thinning high-volume log streams while preserving periodicity.
//
// Example — keep every third line starting at index 1:
//
//	s, _ := stripe.New(stripe.Options{Every: 3, Offset: 1})
//	out := s.Apply(lines)
package stripe
