package watch

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"
)

// Watcher tails a file and emits new lines as they are appended.
type Watcher struct {
	path     string
	interval time.Duration
	output   io.Writer
}

// Options configures a Watcher.
type// Path is the file to watch.
	Path string
	// PollInterval is how often to check for new content.
	PollInterval time.Duration
	// Output is where new lines are written.
	Output io.Writer
}

// NewWatcher creates a Watcher from the given options.
func NewWatcher(opts Options) (*Watcher, error) {
	if opts.Path == "" {
		return nil, fmt.Errorf("watch: path must not be empty")
	}
	if opts.Output == nil {
		return nil, fmt.Errorf("watch: output writer must not be nil")
	}
	interval := opts.PollInterval
	if interval <= 0 {
		interval = 500 * time.Millisecond
	}
	return &Watcher{
		path:     opts.Path,
		interval: interval,
		output:   opts.Output,
	}, nil
}

// Run starts watching the file, writing new lines to output until ctx is cancelled.
func (w *Watcher) Run(ctx context.Context) error {
	f, err := os.Open(w.path)
	if err != nil {
		return fmt.Errorf("watch: open %q: %w", w.path, err)
	}
	defer f.Close()

	// Seek to end so we only emit new content.
	if _, err := f.Seek(0, io.SeekEnd); err != nil {
		return fmt.Errorf("watch: seek %q: %w", w.path, err)
	}

	buf := make([]byte, 4096)
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			for {
				n, err := f.Read(buf)
				if n > 0 {
					if _, werr := w.output.Write(buf[:n]); werr != nil {
						return fmt.Errorf("watch: write output: %w", werr)
					}
				}
				if err == io.EOF {
					break
				}
				if err != nil {
					return fmt.Errorf("watch: read %q: %w", w.path, err)
				}
			}
		}
	}
}
