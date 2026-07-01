package zxipv6wry

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/zu1k/nali/internal/mmap"
	"github.com/zu1k/nali/pkg/wry"
)

// ZXwry is a ZX IPv6 (zxinc) geolocation database reader.
type ZXwry struct {
	wry.IPDB[uint64]
}

// NewZXwry opens the ZX IPv6 database at filePath, downloading it if the file is absent.
func NewZXwry(filePath string) (*ZXwry, error) {
	_, err := os.Stat(filePath)
	if err != nil && os.IsNotExist(err) {
		log.Println("文件不存在，尝试从网络获取最新ZX IPv6数据库")
		_, err = Download(context.Background(), filePath)
		if err != nil {
			return nil, err
		}
	}

	fileData, err := mmap.MapFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("mmap zxipv6wry database: %w", err)
	}

	if !CheckFile(fileData) {
		return nil, errors.New("ZX IPv6数据库存在错误，请重新下载")
	}

	header := fileData[:24]
	offLen := header[6]
	ipLen := header[7]

	start := binary.LittleEndian.Uint64(header[16:24])
	counts := binary.LittleEndian.Uint64(header[8:16])
	end := start + counts*11

	return &ZXwry{
		IPDB: wry.IPDB[uint64]{
			Data: fileData,

			OffLen:   offLen,
			IPLen:    ipLen,
			IPCnt:    counts,
			IdxStart: start,
			IdxEnd:   end,
		},
	}, nil
}

// Find looks up query and returns its location, or an error for non-IPv6 input.
func (db *ZXwry) Find(query string) (result fmt.Stringer, err error) {
	ip := net.ParseIP(query)
	if ip == nil {
		return nil, errors.New("query should be IPv6")
	}
	ip6 := ip.To16()
	if ip6 == nil {
		return nil, errors.New("query should be IPv6")
	}
	ip6 = ip6[:8]
	ipu64 := binary.BigEndian.Uint64(ip6)

	offset := db.SearchIndexV6(ipu64)
	reader := wry.NewReader(db.Data)
	reader.Parse(offset)
	if err := reader.Err(); err != nil {
		return nil, err
	}
	return reader.Result, nil
}

// Name returns the database name.
func (db *ZXwry) Name() string {
	return "zxipv6wry"
}

// CheckFile reports whether data looks like a valid ZX IPv6 database.
func CheckFile(data []byte) bool {
	if len(data) < 4 {
		return false
	}
	if string(data[:4]) != "IPDB" {
		return false
	}

	if len(data) < 24 {
		return false
	}
	header := data[:24]
	start := binary.LittleEndian.Uint64(header[16:24])
	counts := binary.LittleEndian.Uint64(header[8:16])
	end := start + counts*11
	if start >= end || uint64(len(data)) < end {
		return false
	}

	return true
}
