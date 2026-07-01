package qqwry

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/zu1k/nali/internal/mmap"
	"github.com/zu1k/nali/pkg/download"
	"github.com/zu1k/nali/pkg/wry"
)

// DownloadUrls are the mirror URLs the QQwry database is fetched from.
var DownloadUrls = []string{
	"https://github.com/metowolf/qqwry.dat/releases/latest/download/qqwry.dat",
	// Other repo:
	// https://github.com/HMBSbige/qqwry // This repository has been archived since Jun 27, 2024.
	// https://github.com/FW27623/qqwry // This repository's dat format will not be maintained after October 2024.
	// https://github.com/metowolf/qqwry.dat
}

// QQwry is a QQwry (纯真) IPv4 geolocation database reader.
type QQwry struct {
	wry.IPDB[uint32]
}

// NewQQwry new database from path
func NewQQwry(filePath string) (*QQwry, error) {
	_, err := os.Stat(filePath)
	if err != nil && os.IsNotExist(err) {
		log.Println("文件不存在，尝试从网络获取最新纯真 IP 库")
		_, err = download.Download(context.Background(), filePath, DownloadUrls...)
		if err != nil {
			return nil, err
		}
	}

	fileData, err := mmap.MapFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("mmap qqwry database: %w", err)
	}

	if !CheckFile(fileData) {
		return nil, errors.New("纯真 IP 库存在错误，请重新下载")
	}

	header := fileData[0:8]
	start := binary.LittleEndian.Uint32(header[:4])
	end := binary.LittleEndian.Uint32(header[4:])

	return &QQwry{
		IPDB: wry.IPDB[uint32]{
			Data: fileData,

			OffLen:   3,
			IPLen:    4,
			IPCnt:    (end-start)/7 + 1,
			IdxStart: start,
			IdxEnd:   end,
		},
	}, nil
}

// Find looks up query and returns its location, or an error for non-IPv4 input.
func (db QQwry) Find(query string) (result fmt.Stringer, err error) {
	ip := net.ParseIP(query)
	if ip == nil {
		return nil, errors.New("query should be IPv4")
	}
	ip4 := ip.To4()
	if ip4 == nil {
		return nil, errors.New("query should be IPv4")
	}
	ip4uint := binary.BigEndian.Uint32(ip4)

	offset := db.SearchIndexV4(ip4uint)
	if offset <= 0 {
		return nil, errors.New("query not valid")
	}

	reader := wry.NewReader(db.Data)
	reader.Parse(offset + 4)
	if err := reader.Err(); err != nil {
		return nil, err
	}
	return reader.Result.DecodeGBK(), nil
}

// Name returns the database name.
func (db QQwry) Name() string {
	return "qqwry"
}

// CheckFile reports whether data looks like a valid QQwry database.
func CheckFile(data []byte) bool {
	if len(data) < 8 {
		return false
	}

	header := data[0:8]
	start := binary.LittleEndian.Uint32(header[:4])
	end := binary.LittleEndian.Uint32(header[4:])

	// Compare in uint64 so neither the len conversion nor end+7 can overflow.
	if start >= end || uint64(len(data)) < uint64(end)+7 {
		return false
	}

	return true
}
