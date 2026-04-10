// Package diff provides utilities for comparing two sequences of log lines
// and producing structured change sets. It supports both unified and side-by-side
// output formats, and integrates with the output formatter for colorized display.
//
// Typical usage:
//
//	changes := diff.Diff(beforeLines, afterLines)
//	for _, c := range changes {
//		fmt.Println(diff.FormatChange(c))
//	}
package diff
