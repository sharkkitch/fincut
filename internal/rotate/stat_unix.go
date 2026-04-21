//go:build !windows

package rotate

import (
	"os"
	"syscall"
)

func inode(fi os.FileInfo) uint64 {
	if sys, ok := fi.Sys().(*syscall.Stat_t); ok {
		return sys.Ino
	}
	return 0
}

func statFile(path string) (os.FileInfo, error) {
	return os.Stat(path)
}
