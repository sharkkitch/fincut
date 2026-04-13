// Package watch provides tail-style file watching for fincut.
//
// It polls a target file at a configurable interval and writes any
// newly appended bytes to an io.Writer. Watching is bounded by a
// context, allowing clean cancellation from the CLI layer.
//
// The watcher seeks to the end of the file on startup, so only bytes
// appended after Watch begins are forwarded to the output writer.
// If the file is truncated or rotated, the watcher resets its offset
// to zero and resumes from the beginning of the new content.
//
// Typical usage:
//
//	w, err := watch.NewWatcher(watch.Options{
//		Path:         "/var/log/app.log",
//		PollInterval: 250 * time.Millisecond,
//		Output:       os.Stdout,
//	})
//	if err != nil { ... }
//	if err := w.Run(ctx); err != nil && !errors.Is(err, context.Canceled) { ... }
package watch
