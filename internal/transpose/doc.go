// Package transpose provides a Transposer that pivots structured log lines
// by rotating rows into columns.
//
// Given N input lines each containing M delimiter-separated fields, Apply
// returns M output lines each containing N fields — the matrix transpose of
// the input.
//
// Example (delimiter=","):
//
//	input:  ["a,b,c", "1,2,3"]
//	output: ["a,1", "b,2", "c,3"]
//
// Jagged input (lines with different field counts) is handled either by
// leaving missing cells empty or by padding them with a configurable fill
// value when PadFields is enabled.
//
// # Usage
//
// Construct a Transposer with the desired options, then call Apply:
//
//	t := transpose.New(
//		transpose.WithDelimiter(","),
//		transpose.WithPadFields("N/A"),
//	)
//	output, err := t.Apply(input)
//
// Apply returns an error if the input is empty or if a line cannot be parsed
// with the configured delimiter.
package transpose
