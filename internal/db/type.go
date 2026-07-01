package db

import (
	"fmt"

	"github.com/zu1k/nali/internal/constant"
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

// DB describes a single database: its name and aliases, on-disk format and
// file location, together with the query types and languages it supports.
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

	filePath := constant.ResolveDBPath(d.File)

	switch d.Format {
	case FormatQQWry:
		db, err = qqwry.NewQQwry(filePath)
	case FormatZXIPv6Wry:
		db, err = zxipv6wry.NewZXwry(filePath)
	case FormatIPIP:
		db, err = ipip.NewIPIP(filePath)
	case FormatMMDB:
		db, err = geoip.NewGeoIP(filePath, selectedLang)
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

// Format identifies the on-disk format of a database file.
type Format string

// Supported database formats.
const (
	FormatMMDB        Format = "mmdb"
	FormatQQWry       Format = "qqwry"
	FormatZXIPv6Wry   Format = "zxipv6wry"
	FormatIPIP        Format = "ipip"
	FormatIP2Region   Format = "ip2region"
	FormatIP2Location Format = "ip2location"

	FormatCDNYml Format = "cdn-yml"
)

// Predefined language selections shared by databases.
var (
	LanguagesAll = []string{"ALL"}
	LanguagesZH  = []string{"zh-CN"}
	LanguagesEN  = []string{"en"}
)

// Type identifies the category of records a database provides.
type Type string

// Supported query types.
const (
	TypeIPv4 Type = "IPv4"
	TypeIPv6 Type = "IPv6"
	TypeCDN  Type = "CDN"
)

// Predefined type selections shared by databases.
var (
	TypesAll  = []Type{TypeIPv4, TypeIPv6, TypeCDN}
	TypesIP   = []Type{TypeIPv4, TypeIPv6}
	TypesIPv4 = []Type{TypeIPv4}
	TypesIPv6 = []Type{TypeIPv6}
	TypesCDN  = []Type{TypeCDN}
)

// List is an ordered collection of databases.
type List []*DB

// NameMap indexes databases by name and alias.
type NameMap map[string]*DB

// TypeMap groups databases by the query type they serve.
type TypeMap map[Type][]*DB

// From populates the map from dbs, indexing each database by its name and aliases.
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

// From populates the map from dbs, grouping databases by each supported type.
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

// Result is a database query result together with the name of its source.
type Result struct {
	Source string
	common.Result
}
