package transpose

import (
	"fmt"
	"strings"
)

// Transposer pivots lines into columns: each input line becomes a column,
// and the Nth token of each line forms the Nth output row.
type Transposer struct {
	delimiter string
	padFields bool
	fillEmpty string
}

// Options configures a Transposer.
type Options struct {
	// Delimiter separates fields within each input line.
	Delimiter string
	// PadFields pads shorter rows with FillEmpty so all rows have equal width.
	PadFields bool
	// FillEmpty is the value used when padding missing fields. Defaults to "".
	FillEmpty string
}

// New creates a Transposer from the given options.
func New(opts Options) (*Transposer, error) {
	if err := validateOptions(opts); err != nil {
		return nil, err
	}
	return &Transposer{
		delimiter: opts.Delimiter,
		padFields: opts.PadFields,
		fillEmpty: opts.FillEmpty,
	}, nil
}

// Apply transposes the input lines and returns the result.
func (t *Transposer) Apply(lines []string) []string {
	if len(lines) == 0 {
		return []string{}
	}

	// Split every line into fields.
	rows := make([][]string, len(lines))
	maxCols := 0
	for i, line := range lines {
		fields := strings.Split(line, t.delimiter)
		rows[i] = fields
		if len(fields) > maxCols {
			maxCols = len(fields)
		}
	}

	// Build transposed output: column j becomes row j.
	out := make([]string, maxCols)
	for col := 0; col < maxCols; col++ {
		parts := make([]string, len(rows))
		for row := 0; row < len(rows); row++ {
			if col < len(rows[row]) {
				parts[row] = rows[row][col]
			} else if t.padFields {
				parts[row] = t.fillEmpty
			} else {
				parts[row] = ""
			}
		}
		out[col] = strings.Join(parts, t.delimiter)
	}
	return out
}

// FormatSummary returns a one-line summary of a transpose operation.
func FormatSummary(inputLines, outputLines int) string {
	return fmt.Sprintf("transposed %d lines → %d lines", inputLines, outputLines)
}
