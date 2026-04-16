package rotate

import (
	"fmt"
	"time"
)

// RotationEvent records metadata about a detected log rotation.
type RotationEvent struct {
	Path      string
	DetectedAt time.Time
	OldInode  uint64
	NewInode  uint64
}

// FormatEvent returns a human-readable summary of a rotation event.
func FormatEvent(e RotationEvent) string {
	return fmt.Sprintf(
		"[rotate] %s rotated at %s (inode %d -> %d)",
		e.Path,
		e.DetectedAt.Format(time.RFC3339),
		e.OldInode,
		e.NewInode,
	)
}

// FormatSummary returns a compact one-line summary suitable for log output.
func FormatSummary(events []RotationEvent) string {
	if len(events) == 0 {
		return "[rotate] no rotations detected"
	}
	return fmt.Sprintf("[rotate] %d rotation(s) detected", len(events))
}
