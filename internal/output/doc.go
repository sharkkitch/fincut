// Package output provides formatting utilities for fincut log output.
//
// It supports three output formats:
//
//   - plain: raw log lines written as-is (default)
//   - json:  each line wrapped in a JSON object: {"line": "..."}
//   - color: ANSI-colored output based on detected log level keywords
//            (ERROR/FATAL=red, WARN=yellow, INFO=green, DEBUG=cyan)
//
// Usage:
//
//	fmt, err := output.NewFormatter(output.FormatColor, os.Stdout)
//	if err != nil {
//		log.Fatal(err)
//	}
//	_ = fmt.WriteLine("INFO: application started")
package output
