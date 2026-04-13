// Package compact removes blank or whitespace-only lines from a line stream,
// optionally collapsing consecutive blank lines into a single blank line.
package compact

import (
	"errors"
	"strings"
)

// Options controls the behaviour of the Compacter.
type Options struct {
	// CollapseBlank replaces runs of blank lines with a single blank line
	// instead of removing all blank lines.
	CollapseBlank bool

	// TrimSpace trims leading and trailing whitespace before deciding whether
	// a line is blank.
	TrimSpace bool
}

// Compacter filters blank lines from a slice of strings.
type Compacter struct {
	opts Options
}

// New returns a new Compacter or an error if the options are invalid.
func New(opts Options) (*Compacter, error) {
	if err := validateOptions(opts); err != nil {
		return nil, err
	}
	return &Compacter{opts: opts}, nil
}

func validateOptions(_ Options) error {
	// No invalid combinations currently; reserved for future constraints.
	return nil
}

// Apply processes lines and returns the compacted result.
func (c *Compacter) Apply(lines []string) ([]string, error) {
	if len(lines) == 0 {
		return []string{}, nil
	}

	out := make([]string, 0, len(lines))
	prevBlank := false

	for _, line := range lines {
		effective := line
		if c.opts.TrimSpace {
			effective = strings.TrimSpace(line)
		}

		isBlank := strings.TrimSpace(effective) == ""

		if isBlank {
			if c.opts.CollapseBlank && !prevBlank {
				out = append(out, "")
			}
			// If CollapseBlank is false we simply skip blank lines.
			prevBlank = true
			continue
		}

		prevBlank = false
		if c.opts.TrimSpace {
			out = append(out, effective)
		} else {
			out = append(out, line)
		}
	}

	_ = errors.New // keep import if needed by future validators
	return out, nil
}
