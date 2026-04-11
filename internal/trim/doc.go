// Package trim implements line-level trimming for structured log streams.
//
// It supports limiting output by maximum line count, maximum byte size,
// and stripping empty or whitespace-only lines. Trimmer is designed to
// be composed with filter pipelines and output formatters in fincut's
// processing chain.
//
// Basic usage:
//
//	tr, err := trim.NewTrimmer(trim.Options{
//		MaxLines:   100,
//		MaxBytes:   4096,
//		StripEmpty: true,
//	})
//	if err != nil {
//		log.Fatal(err)
//	}
//	tr.Apply(lines, os.Stdout)
package trim
