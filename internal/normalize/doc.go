// Package normalize provides line-level normalization for structured log
// processing pipelines.
//
// A Normalizer applies one or more transformations to each line:
//
//   - TrimSpace: strips leading and trailing whitespace
//   - CollapseSpaces: collapses consecutive whitespace into a single space
//   - Lowercase: converts the entire line to lowercase
//   - StripControl: removes non-printable control characters (preserves tab)
//
// At least one option must be enabled. Options are applied in the following
// fixed order: StripControl → CollapseSpaces → TrimSpace → Lowercase.
//
// Example:
//
//	n, err := normalize.New(normalize.Options{
//		TrimSpace:      true,
//		CollapseSpaces: true,
//	})
//	if err != nil {
//		log.Fatal(err)
//	}
//	result := n.Apply(lines)
package normalize
