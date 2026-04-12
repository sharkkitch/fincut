// Package template provides line-level text transformation using Go's
// text/template engine.
//
// Each input line is rendered through a user-supplied template. An optional
// regular expression with named capture groups can be specified; when a line
// matches, the capture group values are exposed as template variables alongside
// the built-in {{.Line}} variable that always holds the raw line text.
//
// Lines that do not match the optional pattern are passed through unchanged,
// making it safe to apply the templater to heterogeneous log streams.
//
// Example usage:
//
//	tr, err := template.New(template.Options{
//	    Template: "[{{.level}}] {{.msg}}",
//	    Pattern:  `(?P<level>\w+)\s+(?P<msg>.+)`,
//	})
//	if err != nil { ... }
//	lines, err := tr.Apply(input)
package template
