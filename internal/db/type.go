package db

import (
	"fmt"

	"github.com/zu1k/nali/pkg/cdn"
	"github.com/zu1k/nali/pkg/common"
	"github.com/zu1k/nali/pkg/dbif"
	"github.com/zu1k/nali/pkg/geoip"
	"github.com/zu1k/nali/pkg/ip2location"
	"github.com/zu1k/nali/pkg/ip2region"
	"github.com/zu1k/nali/pkg/ipip"
	"github.com/zu1k/nali/pkg/qqwry"
	"github.com/zu1k/nali/pkg/zxipv6wry"
)

type DB struct {
	Name      string
	NameAlias []string `yaml:"name-alias,omitempty" mapstructure:"name-alias"`
	Format    Format
	File      string

	Languages []string
	Types     []Type

	DownloadUrls []string `yaml:"download-urls,omitempty" mapstructure:"download-urls"`
}

func (d *DB) get() (db dbif.DB, err error) {
	cacheMu.RLock()
	cached, found := dbNameCache[d.Name]
	cacheMu.RUnlock()
	if found {
		return cached, nil
	}

	filePath := d.File

	switch d.Format {
	case FormatQQWry:
		db, err = qqwry.NewQQwry(filePath)
	case FormatZXIPv6Wry:
		db, err = zxipv6wry.NewZXwry(filePath)
	case FormatIPIP:
		db, err = ipip.NewIPIP(filePath)
	case FormatMMDB:
		db, err = geoip.NewGeoIP(filePath)
	case FormatIP2Region:
		db, err = ip2region.NewIp2Region(filePath)
	case FormatIP2Location:
		db, err = ip2location.NewIP2Location(filePath)
	case FormatCDNYml:
		db, err = cdn.NewCDN(filePath)
	default:
		return nil, fmt.Errorf("DB format not supported: %s", d.Format)
	}

	if err != nil {
		return nil, fmt.Errorf("database init failed: %w", err)
	}

	cacheMu.Lock()
	dbNameCache[d.Name] = db
	cacheMu.Unlock()
	return
}

type Format string

const (
	FormatMMDB        Format = "mmdb"
	FormatQQWry              = "qqwry"
	FormatZXIPv6Wry          = "zxipv6wry"
	FormatIPIP               = "ipip"
	FormatIP2Region          = "ip2region"
	FormatIP2Location        = "ip2location"

	FormatCDNYml = "cdn-yml"
)

var (
	LanguagesAll = []string{"ALL"}
	LanguagesZH  = []string{"zh-CN"}
	LanguagesEN  = []string{"en"}
)

type Type string

const (
	TypeIPv4 Type = "IPv4"
	TypeIPv6      = "IPv6"
	TypeCDN       = "CDN"
)

var (
	TypesAll  = []Type{TypeIPv4, TypeIPv6, TypeCDN}
	TypesIP   = []Type{TypeIPv4, TypeIPv6}
	TypesIPv4 = []Type{TypeIPv4}
	TypesIPv6 = []Type{TypeIPv6}
	TypesCDN  = []Type{TypeCDN}
)

type List []*DB
type NameMap map[string]*DB
type TypeMap map[Type][]*DB

func (m *NameMap) From(dbs List) {
	for _, db := range dbs {
		(*m)[db.Name] = db

		if alias := db.NameAlias; alias != nil {
			for _, aName := range alias {
				(*m)[aName] = db
			}
		}
	}
}

func (m *TypeMap) From(dbs List) {
	for _, db := range dbs {
		for _, typ := range db.Types {
			dbsInType := (*m)[typ]
			if dbsInType == nil {
				dbsInType = []*DB{db}
			} else {
				dbsInType = append(dbsInType, db)
			}
			(*m)[typ] = dbsInType
		}
	}
}

func getDbByName(name string) (*DB, error) {
	if dbInfo, found := NameDBMap[name]; found {
		return dbInfo, nil
	}

	defaultNameDBMap := NameMap{}
	defaultNameDBMap.From(GetDefaultDBList())
	if dbInfo, found := defaultNameDBMap[name]; found {
		return dbInfo, nil
	}

	return nil, fmt.Errorf("DB with name %s not found", name)
}

type Result struct {
	Source string
	common.Result
}
