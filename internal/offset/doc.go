// Package offset computes per-line byte offsets within a slice of log lines.
//
// It is useful when you need to map a line number back to an exact byte
// position in the original file — for example, when building a seekable
// index or reporting precise locations to the user.
//
// Basic usage:
//
//	o, err := offset.New(offset.Options{StartLine: 1, EndLine: 50})
//	if err != nil {
//		log.Fatal(err)
//	}
//	entries := o.Apply(lines)
//	for _, e := range entries {
//		fmt.Printf("line %d  bytes %d–%d  %s\n", e.Line, e.ByteStart, e.ByteEnd, e.Content)
//	}
package offset
