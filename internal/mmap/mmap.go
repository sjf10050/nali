// Package mmap provides a minimal cross-platform read-only memory-mapped file.
// Short‑lived CLI tool: callers are not required to call Unmap — the OS
// reclaims the mapping on process exit. Unmap is provided for completeness.
package mmap

import (
	"errors"
	"os"
)

// ErrEmptyFile is returned when trying to map a zero‑byte file.
var ErrEmptyFile = errors.New("mmap: file is empty")

// MapFile maps the entire file at path into memory for read‑only access.
// The returned slice is valid for the lifetime of the mapping (typically
// the process lifetime for a CLI tool).
func MapFile(path string) ([]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		return nil, err
	}
	if fi.Size() == 0 {
		return nil, ErrEmptyFile
	}

	return mapFile(f, fi.Size())
}

// Unmap releases the memory mapping. After this call the slice must not be
// accessed. Not required for short‑lived processes — the OS cleans up on exit.
func Unmap(data []byte) error {
	return unmapFile(data)
}
