// Package summarize counts distinct lines across a stream and reports
// frequency summaries.
//
// Basic usage:
//
//	s, err := summarize.New(summarize.Options{TopN: 10, CaseSensitive: false})
//	if err != nil { ... }
//	for _, line := range lines {
//		s.Add(line)
//	}
//	results := s.Results()
//	fmt.Print(summarize.FormatSummary(results, s.Total()))
//
// Options:
//   - TopN          – keep only the N most frequent lines (0 = unlimited).
//   - MinCount      – omit lines seen fewer than MinCount times (0/1 = all).
//   - CaseSensitive – when false, "ERROR" and "error" are counted together.
package summarize
