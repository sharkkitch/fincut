// Package prefix provides a Prefixer that prepends a fixed label or
// sequential line numbers to each line of input.
//
// Usage:
//
//	p, err := prefix.New(prefix.Options{
//		LineNumbers: true,
//		Width:       4,   // zero-pad to 4 digits
//		Separator:   " | ",
//	})
//	if err != nil {
//		log.Fatal(err)
//	}
//	result := p.Apply(lines)
//
// When Text is set, every line receives the same static prefix.
// When LineNumbers is true, each line is prefixed with its 1-based index.
// The two modes are mutually exclusive.
package prefix
