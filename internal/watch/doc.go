// Package watch provides tail-style file watching for fincut.
//
// It polls a target file at a configurable interval and writes any
// newly appended bytes to an io.Writer. Watching is bounded by a
// context, allowing clean cancellation from the CLI layer.
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
