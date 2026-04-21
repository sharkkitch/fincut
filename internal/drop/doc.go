// Package drop provides a line-dropping filter that removes lines from a
// slice based on one or more regular-expression patterns.
//
// Basic usage:
//
//	d, err := drop.New(drop.Options{
//		Patterns: []string{`^DEBUG`, `^TRACE`},
//	})
//	if err != nil {
//		log.Fatal(err)
//	}
//	filtered := d.Apply(lines)
//
// Setting Invert to true keeps only the lines that match at least one
// pattern, effectively turning the dropper into an inclusive filter.
package drop
