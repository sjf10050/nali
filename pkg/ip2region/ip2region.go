package ip2region

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/zu1k/nali/pkg/download"
	"github.com/zu1k/nali/pkg/wry"

	"github.com/lionsoul2014/ip2region/binding/golang/xdb"
)

// DownloadUrls are the mirror URLs the ip2region database is fetched from.
var DownloadUrls = []string{
	"https://cdn.jsdelivr.net/gh/lionsoul2014/ip2region/data/ip2region.xdb",
	"https://raw.githubusercontent.com/lionsoul2014/ip2region/master/data/ip2region.xdb",
}

// Ip2Region is an ip2region xdb database reader.
type Ip2Region struct {
	seacher *xdb.Searcher
}

// NewIp2Region opens the ip2region xdb at filePath, downloading it if the file is absent.
func NewIp2Region(filePath string) (*Ip2Region, error) {
	_, err := os.Stat(filePath)
	if err != nil && os.IsNotExist(err) {
		log.Println("文件不存在，尝试从网络获取最新 ip2region 库")
		_, err = download.Download(context.Background(), filePath, DownloadUrls...)
		if err != nil {
			return nil, err
		}
	}

	// Use memory-mapped file access instead of loading entire file into memory
	searcher, err := xdb.NewWithFileOnly(xdb.IPv4, filePath)
	if err != nil {
		fmt.Printf("无法解析 ip2region xdb 数据库: %s\n", err)
		return nil, err
	}
	return &Ip2Region{
		seacher: searcher,
	}, nil
}

// Find looks up query and returns its region information.
func (db Ip2Region) Find(query string) (result fmt.Stringer, err error) {
	if db.seacher != nil {
		res, err := db.seacher.Search(query)
		if err != nil {
			return nil, err
		} else {
			return wry.Result{
				Country: strings.ReplaceAll(res, "|0", ""),
			}, nil
		}
	}

	return nil, errors.New("ip2region 未初始化")
}

// Name returns the database name.
func (db Ip2Region) Name() string {
	return "ip2region"
}
