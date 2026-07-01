//go:build windows

package mmap

import (
	"os"
	"syscall"
	"unsafe"
)

func mapFile(f *os.File, size int64) ([]byte, error) {
	h, err := syscall.CreateFileMapping(syscall.Handle(f.Fd()), nil, syscall.PAGE_READONLY, 0, 0, nil)
	if err != nil {
		return nil, os.NewSyscallError("CreateFileMapping", err)
	}
	defer syscall.CloseHandle(h)

	addr, err := syscall.MapViewOfFile(h, syscall.FILE_MAP_READ, 0, 0, 0)
	if err != nil {
		return nil, os.NewSyscallError("MapViewOfFile", err)
	}

	// The mapping handle has been closed; the view keeps a reference.
	return unsafe.Slice((*byte)(unsafe.Pointer(addr)), int(size)), nil
}

func unmapFile(data []byte) error {
	if len(data) == 0 {
		return nil
	}
	return syscall.UnmapViewOfFile(uintptr(unsafe.Pointer(&data[0])))
}
