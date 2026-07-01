//go:build !windows

package mmap

import (
	"os"

	"golang.org/x/sys/unix"
)

func mapFile(f *os.File, size int64) ([]byte, error) {
	return unix.Mmap(int(f.Fd()), 0, int(size), unix.PROT_READ, unix.MAP_SHARED)
}

func unmapFile(data []byte) error {
	return unix.Munmap(data)
}
