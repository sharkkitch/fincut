package stats

import (
	"fmt"
	"io"
)

// Format controls how a report is rendered.
type Format string

const (
	FormatText Format = "text"
	FormatJSON Format = "json"
)

// Reporter writes a Collector summary to an io.Writer.
type Reporter struct {
	w      io.Writer
	format Format
}

// NewReporter creates a Reporter writing to w in the given format.
// Returns an error if the format is unsupported.
func NewReporter(w io.Writer, format Format) (*Reporter, error) {
	switch format {
	case FormatText, FormatJSON:
		// valid
	default:
		return nil, fmt.Errorf("stats: unsupported format %q", format)
	}
	return &Reporter{w: w, format: format}, nil
}

// Write renders the collector's statistics to the reporter's writer.
func (r *Reporter) Write(c *Collector) error {
	var output string
	switch r.format {
	case FormatJSON:
		output = fmt.Sprintf(
			`{"lines_read":%d,"lines_matched":%d,"lines_dropped":%d,`+
				`"bytes_read":%d,"filter_stages":%d,"match_rate":%.4f,"elapsed_ms":%d}\n`,
			c.LinesRead,
			c.LinesMatched,
			c.LinesDropped,
			c.BytesRead,
			c.FilterStages,
			c.MatchRate(),
			c.Elapsed().Milliseconds(),
		)
	default:
		output = c.Summary() + "\n"
	}
	_, err := fmt.Fprint(r.w, output)
	return err
}
