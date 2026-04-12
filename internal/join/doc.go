// Package join provides a Joiner that merges consecutive lines into
// single records using a configurable separator and group size.
//
// # Basic usage
//
//	j, err := join.New(join.Options{
//		GroupSize: 3,
//		Separator: " | ",
//	})
//	if err != nil {
//		log.Fatal(err)
//	}
//	result := j.Apply(lines)
//
// A GroupSize of 0 joins all input lines into a single output line.
package join
