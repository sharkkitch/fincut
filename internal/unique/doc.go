// Package unique implements a line deduplication filter that retains only the
// first occurrence of each unique key.
//
// The uniqueness key can be:
//   - the entire line (default, Field == 0)
//   - a specific delimited field (Field > 0, Delimiter required)
//
// Comparison can optionally be made case-insensitive.
//
// Example usage:
//
//	u, err := unique.New(unique.Options{
//		Field:           2,
//		Delimiter:       ",",
//		CaseInsensitive: true,
//	})
//	if err != nil {
//		log.Fatal(err)
//	}
//	result := u.Apply(lines)
package unique
