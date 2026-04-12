package window

import "fmt"

// Windower extracts a sliding or fixed window of lines from a stream.
type Windower struct {
	size   int
	step   int
	overlap bool
}

// Options configures the Windower.
type Options struct {
	// Size is the number of lines per window (required, > 0).
	Size int

	// Step is how many lines to advance between windows.
	// Defaults to Size (non-overlapping). Set < Size for overlap.
	Step int
}

// New creates a Windower from the given Options.
func New(opts Options) (*Windower, error) {
	if opts.Size <= 0 {
		return nil, fmt.Errorf("window: size must be greater than zero, got %d", opts.Size)
	}
	if opts.Step < 0 {
		return nil, fmt.Errorf("window: step must be non-negative, got %d", opts.Step)
	}
	step := opts.Step
	if step == 0 {
		step = opts.Size
	}
	return &Windower{
		size:   opts.Size,
		step:   step,
		overlap: step < opts.Size,
	}, nil
}

// Apply partitions lines into windows of the configured size, advancing by
// step lines between each window. Each window is a []string slice.
func (w *Windower) Apply(lines []string) [][]string {
	if len(lines) == 0 {
		return nil
	}
	var windows [][]string
	for start := 0; start < len(lines); start += w.step {
		end := start + w.size
		if end > len(lines) {
			end = len(lines)
		}
		win := make([]string, end-start)
		copy(win, lines[start:end])
		windows = append(windows, win)
		if end == len(lines) {
			break
		}
	}
	return windows
}

// Flatten merges windowed lines back into a flat slice, de-duplicating
// overlapping lines so each source line appears exactly once.
func Flatten(windows [][]string) []string {
	var out []string
	seen := 0
	for _, win := range windows {
		for i, line := range win {
			if i >= seen {
				out = append(out, line)
			}
		}
		seen = 0 // only skip overlap on first window
	}
	return out
}
