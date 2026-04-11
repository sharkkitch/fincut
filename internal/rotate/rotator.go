// Package rotate provides log rotation detection and handling for fincut.
// It detects when a log file has been rotated (truncated or replaced) and
// allows consumers to reset their read position accordingly.
package rotate

import (
	"errors"
	"io"
	"os"
	"time"
)

// Options configures the Rotator behaviour.
type Options struct {
	Path     string
	Interval time.Duration
	Output   io.Writer
}

// Rotator watches a file for rotation events.
type Rotator struct {
	path     string
	interval time.Duration
	output   io.Writer
	lastSize int64
	lastIno  uint64
}

// NewRotator creates a new Rotator with the given options.
func NewRotator(opts Options) (*Rotator, error) {
	if opts.Path == "" {
		return nil, errors.New("rotate: path must not be empty")
	}
	if opts.Output == nil {
		return nil, errors.New("rotate: output writer must not be nil")
	}
	if opts.Interval <= 0 {
		opts.Interval = 2 * time.Second
	}
	return &Rotator{
		path:     opts.Path,
		interval: opts.Interval,
		output:   opts.Output,
	}, nil
}

// Snapshot captures the current inode and size of the watched file.
func (r *Rotator) Snapshot() error {
	info, err := os.Stat(r.path)
	if err != nil {
		return err
	}
	r.lastSize = info.Size()
	r.lastIno = inode(info)
	return nil
}

// Detect returns true if the file has been rotated since the last Snapshot.
func (r *Rotator) Detect() (bool, error) {
	info, err := os.Stat(r.path)
	if err != nil {
		return false, err
	}
	currentIno := inode(info)
	currentSize := info.Size()
	if currentIno != r.lastIno || currentSize < r.lastSize {
		return true, nil
	}
	return false, nil
}

// Interval returns the polling interval.
func (r *Rotator) Interval() time.Duration {
	return r.interval
}
