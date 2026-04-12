package pivot

import (
	"fmt"
	"sort"
	"strings"
)

// Pivotter groups lines by a key field and aggregates values.
type Pivotter struct {
	keyField   int
	valueField int
	delimiter  string
	aggregator string
	groups     map[string][]string
}

// New creates a Pivotter with the given options.
func New(opts Options) (*Pivotter, error) {
	if err := validateOptions(opts); err != nil {
		return nil, err
	}
	agg := opts.Aggregator
	if agg == "" {
		agg = "count"
	}
	return &Pivotter{
		keyField:   opts.KeyField,
		valueField: opts.ValueField,
		delimiter:  opts.Delimiter,
		aggregator: agg,
		groups:     make(map[string][]string),
	}, nil
}

// Apply processes lines and returns pivoted output rows.
func (p *Pivotter) Apply(lines []string) []string {
	p.groups = make(map[string][]string)
	for _, line := range lines {
		parts := strings.Split(line, p.delimiter)
		if p.keyField >= len(parts) {
			continue
		}
		key := strings.TrimSpace(parts[p.keyField])
		var val string
		if p.valueField >= 0 && p.valueField < len(parts) {
			val = strings.TrimSpace(parts[p.valueField])
		}
		p.groups[key] = append(p.groups[key], val)
	}
	return p.format()
}

func (p *Pivotter) format() []string {
	keys := make([]string, 0, len(p.groups))
	for k := range p.groups {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	result := make([]string, 0, len(keys))
	for _, k := range keys {
		vals := p.groups[k]
		var agg string
		switch p.aggregator {
		case "count":
			agg = fmt.Sprintf("%d", len(vals))
		case "sum":
			sum := 0.0
			for _, v := range vals {
				var n float64
				fmt.Sscanf(v, "%f", &n)
				sum += n
			}
			agg = fmt.Sprintf("%g", sum)
		case "values":
			agg = strings.Join(vals, ",")
		default:
			agg = fmt.Sprintf("%d", len(vals))
		}
		result = append(result, fmt.Sprintf("%s%s%s", k, p.delimiter, agg))
	}
	return result
}
