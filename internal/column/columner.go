package column

import (
	"fmt"
	"strings"
)

// Columner extracts and formats specific columns from delimited text lines.
type Columner struct {
	delimiter string
	fields    []int
	outputSep string
	strict    bool
}

// Options configures Columner behaviour.
type Options struct {
	// Delimiter is the field separator (required).
	Delimiter string
	// Fields is the 1-based list of column indices to extract (required).
	Fields []int
	// OutputSep is the separator used when joining extracted fields (default: same as Delimiter).
	OutputSep string
	// Strict causes Apply to return an error when a line has fewer fields than requested.
	Strict bool
}

// New creates a Columner from opts.
func New(opts Options) (*Columner, error) {
	if err := validateOptions(opts); err != nil {
		return nil, err
	}
	sep := opts.OutputSep
	if sep == "" {
		sep = opts.Delimiter
	}
	return &Columner{
		delimiter: opts.Delimiter,
		fields:    opts.Fields,
		outputSep: sep,
		strict:    opts.Strict,
	}, nil
}

// Apply processes lines and returns only the requested columns.
func (c *Columner) Apply(lines []string) ([]string, error) {
	out := make([]string, 0, len(lines))
	for i, line := range lines {
		parts := strings.Split(line, c.delimiter)
		selected := make([]string, 0, len(c.fields))
		for _, f := range c.fields {
			idx := f - 1
			if idx < 0 || idx >= len(parts) {
				if c.strict {
					return nil, fmt.Errorf("line %d: field %d out of range (have %d fields)", i+1, f, len(parts))
				}
				selected = append(selected, "")
				continue
			}
			selected = append(selected, parts[idx])
		}
		out = append(out, strings.Join(selected, c.outputSep))
	}
	return out, nil
}
