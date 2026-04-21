package rotate

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

// Rotator monitors a file for rotation events (truncation or inode change)
// and notifies via a channel when rotation is detected.
type Rotator struct {
	path     string
	out      io.Writer
	interval time.Duration
	mu       sync.Mutex
	lastSize int64
	lastIno  uint64
}

// Options configures a Rotator.
type Options struct {
	Path     string
	Out      io.Writer
	Interval time.Duration
}

// NewRotator creates a new Rotator with the given options.
func NewRotator(opts Options) (*Rotator, error) {
	if err := validateOptions(opts); err != nil {
		return nil, err
	}
	if opts.Interval == 0 {
		opts.Interval = time.Second
	}
	r := &Rotator{
		path:     opts.Path,
		out:      opts.Out,
		interval: opts.Interval,
	}
	if err := r.snapshot(); err != nil {
		return nil, fmt.Errorf("rotate: initial snapshot: %w", err)
	}
	return r, nil
}

func (r *Rotator) snapshot() error {
	fi, err := os.Stat(r.path)
	if err != nil {
		return err
	}
	r.lastSize = fi.Size()
	r.lastIno = inode(fi)
	return nil
}

// Detect checks whether the file has been rotated since the last call.
// It returns true and emits a formatted event to Out when rotation is detected.
func (r *Rotator) Detect() (bool, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	fi, err := os.Stat(r.path)
	if err != nil {
		return false, fmt.Errorf("rotate: stat: %w", err)
	}

	curSize := fi.Size()
	curIno := inode(fi)

	rotated := curIno != r.lastIno || curSize < r.lastSize
	if rotated {
		event := FormatEvent(r.path, r.lastIno, curIno, r.lastSize, curSize)
		fmt.Fprintln(r.out, event)
		r.lastSize = curSize
		r.lastIno = curIno
	}
	return rotated, nil
}
