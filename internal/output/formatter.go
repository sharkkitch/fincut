package output

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

// Format represents the output format for log lines.
type Format string

const (
	FormatPlain Format = "plain"
	FormatJSON  Format = "json"
	FormatColor Format = "color"
)

// Formatter writes log lines to an output writer in a specified format.
type Formatter struct {
	format Format
	w      io.Writer
}

// NewFormatter creates a new Formatter with the given format and writer.
func NewFormatter(format Format, w io.Writer) (*Formatter, error) {
	switch format {
	case FormatPlain, FormatJSON, FormatColor:
		return &Formatter{format: format, w: w}, nil
	default:
		return nil, fmt.Errorf("unsupported format: %q", format)
	}
}

// WriteLine writes a single log line according to the configured format.
func (f *Formatter) WriteLine(line string) error {
	var out string
	switch f.format {
	case FormatJSON:
		b, err := json.Marshal(map[string]string{"line": line})
		if err != nil {
			return fmt.Errorf("json marshal: %w", err)
		}
		out = string(b)
	case FormatColor:
		out = colorize(line)
	default:
		out = line
	}
	_, err := fmt.Fprintln(f.w, out)
	return err
}

// colorize applies a simple ANSI color to lines containing common log levels.
func colorize(line string) string {
	upper := strings.ToUpper(line)
	switch {
	case strings.Contains(upper, "ERROR") || strings.Contains(upper, "FATAL"):
		return "\033[31m" + line + "\033[0m" // red
	case strings.Contains(upper, "WARN"):
		return "\033[33m" + line + "\033[0m" // yellow
	case strings.Contains(upper, "INFO"):
		return "\033[32m" + line + "\033[0m" // green
	case strings.Contains(upper, "DEBUG"):
		return "\033[36m" + line + "\033[0m" // cyan
	default:
		return line
	}
}
