// Package annotate provides the Annotator type, which prepends structured
// metadata to each line of a log stream.
//
// Supported annotation modes:
//   - Line numbers (1-based)
//   - UTC timestamps with a configurable Go time-format string
//   - A static prefix string
//
// Multiple modes may be combined; fields are joined by a configurable
// separator (default " | ").
//
// Example:
//
//	a, _ := annotate.New(annotate.Options{LineNumbers: true, Prefix: "svc"})
//	out := a.Apply([]string{"hello", "world"})
//	// out[0] == "1 | svc | hello"
package annotate
