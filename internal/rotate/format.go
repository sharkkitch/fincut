package rotate

import (
	"fmt"
	"time"
)

// FormatEvent returns a human-readable rotation event string.
func FormatEvent(path string, oldIno, newIno uint64, oldSize, newSize int64) string {
	ts := time.Now().UTC().Format(time.RFC3339)
	if oldIno != newIno {
		return fmt.Sprintf(
			"[%s] rotate: inode changed path=%s old_inode=%d new_inode=%d",
			ts, path, oldIno, newIno,
		)
	}
	return fmt.Sprintf(
		"[%s] rotate: truncation detected path=%s old_size=%d new_size=%d",
		ts, path, oldSize, newSize,
	)
}

// FormatSummary returns a brief summary string for reporting.
func FormatSummary(detections int, path string) string {
	return fmt.Sprintf("rotate: %d rotation(s) detected for %s", detections, path)
}
