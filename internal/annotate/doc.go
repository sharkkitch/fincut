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
// # Options
//
// Options.TimestampFormat accepts any layout string recognised by the
// standard library's [time.Format]. Common layouts include [time.RFC3339]
// and [time.DateTime]. When TimestampFormat is empty, timestamps are
// omitted even if timestamp annotation is otherwise enabled.
//
// # Example
//
//	a, _ := annotate.New(annotate.Options{LineNumbers: true, Prefix: "svc"})
//	out := a.Apply([]string{"hello", "world"})
//	// out[0] == "1 | svc | hello"
//	// out[1] == "2 | svc | world"
package annotate
