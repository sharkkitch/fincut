//go:build !windows

package rotate

import (
	"os"
	"syscall"
)

// inode returns the inode number of the given FileInfo on Unix systems.
func inode(fi os.FileInfo) uint64 {
	if fi == nil {
		return 0
	}
	if stat, ok := fi.Sys().(*syscall.Stat_t); ok {
		return stat.Ino
	}
	return 0
}
