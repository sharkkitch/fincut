// Package freq implements line-frequency analysis for fincut pipelines.
//
// It counts how often each unique value appears across a slice of log lines,
// supporting optional field extraction (by whitespace-delimited position),
// case-insensitive matching, and top-N result capping.
//
// Basic usage:
//
//	c, err := freq.New(freq.Options{TopN: 10, Field: 2})
//	if err != nil {
//		log.Fatal(err)
//	}
//	c.Add(lines)
//	for _, e := range c.Results() {
//		fmt.Printf("%6d  %s\n", e.Count, e.Value)
//	}
package freq
