// Package dedup provides line-level deduplication for log streams.
//
// A Deduper tracks previously seen lines using SHA-256 hashes and filters
// out duplicates on subsequent calls to Apply. It supports an optional
// sliding window to bound memory usage, and case-insensitive comparison
// for normalised log output.
//
// Basic usage:
//
//	d, err := dedup.New(dedup.Options{
//		WindowSize:    500,
//		CaseSensitive: false,
//	})
//	if err != nil {
//		log.Fatal(err)
//	}
//	unique := d.Apply(lines)
package dedup
