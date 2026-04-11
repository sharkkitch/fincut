// Package mask provides a Masker that redacts sensitive patterns from log lines
// using regular expressions.
//
// A Masker is constructed with one or more regex patterns and an optional
// replacement string. When applied to a slice of log lines, every substring
// matching any pattern is replaced with the configured replacement text
// (defaulting to "[REDACTED]").
//
// Typical usage:
//
//	m, err := mask.New(mask.Options{
//		Patterns:    []string{`password=\S+`, `api_key=[A-Za-z0-9]+`},
//		Replacement: "[REDACTED]",
//	})
//	if err != nil {
//		log.Fatal(err)
//	}
//	clean := m.Apply(lines)
package mask
