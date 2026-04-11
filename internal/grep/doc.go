// Package grep provides pattern-based line searching with optional context
// windows for the fincut CLI tool.
//
// A Grepper is constructed with one or more regular expression patterns and
// optional before/after context line counts. It scans a slice of log lines and
// returns Match values that include the matched line, its 1-based line number,
// and the surrounding context lines.
//
// Invert mode returns lines that do NOT match any of the supplied patterns,
// mirroring the behaviour of grep -v.
//
// Example usage:
//
//	g, err := grep.New(grep.Options{
//		Patterns:      []string{"ERROR", "FATAL"},
//		ContextBefore: 2,
//		ContextAfter:  1,
//	})
//	matches := g.Apply(lines)
package grep
