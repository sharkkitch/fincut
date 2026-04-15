// Package align implements column alignment for structured log lines.
//
// It splits each line on a configurable delimiter, measures the maximum
// field width per column across all lines, then re-joins the fields with
// uniform padding so columns line up visually.
//
// Example:
//
//	a, _ := align.New(align.Options{Delimiter: "|", Padding: 2})
//	out := a.Apply(lines)
//
// Fields are left-padded to the widest value in their column; the final
// field on each line is never padded so trailing whitespace is avoided.
package align
