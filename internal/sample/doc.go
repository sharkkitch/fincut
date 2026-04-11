// Package sample implements line-level sampling strategies for fincut log pipelines.
//
// Two mutually exclusive sampling modes are supported:
//
//   - Rate-based: keeps every Nth line (e.g. Rate=5 retains lines 0, 5, 10, …).
//     Useful for deterministic, evenly-spaced thinning of high-volume logs.
//
//   - Probability-based: retains each line independently with probability P ∈ (0,1].
//     A fixed Seed may be provided for reproducible output.
//
// Example usage:
//
//	s, err := sample.New(sample.Options{Rate: 10})
//	if err != nil {
//		log.Fatal(err)
//	}
//	sampled := s.Apply(lines)
package sample
