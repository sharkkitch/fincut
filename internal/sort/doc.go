// Package sort provides line-level sorting for fincut pipelines.
//
// A Sorter can order lines lexicographically (ascending or descending),
// optionally deduplicate adjacent equal keys, and extract a sort key from
// a delimited field rather than the whole line.
//
// Example usage:
//
//	s, err := sort.New(sort.Options{
//		Reverse:   false,
//		Unique:    true,
//		Delimiter: "|",
//		Field:     2,
//	})
//	if err != nil {
//		log.Fatal(err)
//	}
//	result := s.Apply(lines)
package sort
