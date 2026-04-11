// Package highlight provides regex-based term highlighting for structured
// log lines. It supports multiple simultaneous patterns, each assigned a
// distinct ANSI color from a rotating palette. An optional bold mode
// emphasises matched terms further.
//
// Usage:
//
//	h, err := highlight.New(highlight.Options{
//		Patterns: []string{"error", "warn"},
//		Bold:     true,
//	})
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Println(h.Apply(line))
//
// Use StripANSI to recover the original text from a highlighted string.
package highlight
