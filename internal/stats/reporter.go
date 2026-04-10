package stats

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
)

// Reporter writes collected statistics to an output writer
// in either plain text or JSON format.
type Reporter struct {
	w      io.Writer
	format string
}

// NewReporter creates a Reporter that writes to w using the given format.
// Supported formats: "plain", "json".
func NewReporter(w io.Writer, format string) (*Reporter, error) {
	if w == nil {
		return nil, errors.New("stats: reporter writer must not be nil")
	}
	switch format {
	case "plain", "json":
		// valid
	default:
		return nil, fmt.Errorf("stats: unsupported report format %q", format)
	}
	return &Reporter{w: w, format: format}, nil
}

// Report writes a summary of the collector's statistics to the reporter's writer.
func (r *Reporter) Report(c *Collector) {
	switch r.format {
	case "json":
		r.writeJSON(c)
	default:
		r.writePlain(c)
	}
}

func (r *Reporter) writePlain(c *Collector) {
	fmt.Fprintf(r.w, "matched:    %d\n", c.Matched())
	fmt.Fprintf(r.w, "dropped:    %d\n", c.Dropped())
	fmt.Fprintf(r.w, "total:      %d\n", c.Total())
	fmt.Fprintf(r.w, "match_rate: %.2f%%\n", c.MatchRate()*100)
	fmt.Fprintf(r.w, "elapsed:    %s\n", c.Elapsed())
}

func (r *Reporter) writeJSON(c *Collector) {
	payload := map[string]interface{}{
		"matched":    c.Matched(),
		"dropped":    c.Dropped(),
		"total":      c.Total(),
		"match_rate": c.MatchRate(),
		"elapsed_ms": c.Elapsed().Milliseconds(),
	}
	data, err := json.Marshal(payload)
	if err != nil {
		fmt.Fprintf(r.w, `{"error":%q}\n`, err.Error())
		return
	}
	r.w.Write(data)
	fmt.Fprintln(r.w)
}
