// Package rotate detects log file rotation by monitoring inode changes
// or file truncation, and re-opens the file when rotation is detected.
package rotate

import (
	"fmt"
	"io"
	"os"
	"time"
)

// Options configures the Rotator.
type Options struct {
	// Path is the log file to monitor.
	Path string
	// Output is where detected rotation events are written.
	Output io.Writer
	// Interval is how often to poll for rotation. Defaults to 5s.
	Interval time.Duration
}

// Rotator monitors a file for rotation events.
type Rotator struct {
	opts Options
}

// NewRotator creates a Rotator with the given options.
func NewRotator(opts Options) (*Rotator, error) {
	if opts.Path == "" {
		return nil, fmt.Errorf("rotate: path must not be empty")
	}
	if opts.Output == nil {
		return nil, fmt.Errorf("rotate: output writer must not be nil")
	}
	if opts.Interval <= 0 {
		opts.Interval = 5 * time.Second
	}
	return &Rotator{opts: opts}, nil
}

// Detect checks whether the file at opts.Path has been rotated since
// the provided baseline os.FileInfo. It returns true if rotation is
// detected (inode changed or file is smaller than before).
func (r *Rotator) Detect(baseline os.FileInfo) (bool, error) {
	current, err := os.Stat(r.opts.Path)
	if err != nil {
		if os.IsNotExist(err) {
			return true, nil
		}
		return false, fmt.Errorf("rotate: stat %s: %w", r.opts.Path, err)
	}
	if baseline == nil {
		return false, nil
	}
	if inode(current) != inode(baseline) {
		return true, nil
	}
	if current.Size() < baseline.Size() {
		return true, nil
	}
	return false, nil
}

// Baseline returns the current os.FileInfo for the monitored file.
func (r *Rotator) Baseline() (os.FileInfo, error) {
	fi, err := os.Stat(r.opts.Path)
	if err != nil {
		return nil, fmt.Errorf("rotate: baseline stat %s: %w", r.opts.Path, err)
	}
	return fi, nil
}
