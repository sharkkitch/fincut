// Package clamp provides numeric range filtering for structured log lines.
//
// It compiles a regex pattern with a capture group that is expected to match
// a numeric value. Lines whose captured value falls within [Min, Max] are
// retained; all others are dropped.
//
// Example usage:
//
//	min := 200.0
//	max := 399.0
//	c, err := clamp.New(clamp.Options{
//		Pattern: `status=(\d+)`,
//		Min:     &min,
//		Max:     &max,
//	})
//	if err != nil {
//		log.Fatal(err)
//	}
//	result := c.Apply(lines)
package clamp
