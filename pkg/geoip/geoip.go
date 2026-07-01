package geoip

import (
	"errors"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/oschwald/geoip2-golang"
)

// GeoIP is a MaxMind GeoIP2/GeoLite2 database reader.
type GeoIP struct {
	db   *geoip2.Reader
	lang string
}

// NewGeoIP opens the GeoIP2 database at filePath, using lang for localized names.
func NewGeoIP(filePath string, lang string) (*GeoIP, error) {
	_, err := os.Stat(filePath)
	if err != nil && os.IsNotExist(err) {
		log.Println("文件不存在，请自行下载 Geoip2 City库，并保存在", filePath)
		return nil, err
	} else {
		db, err := geoip2.Open(filePath)
		if err != nil {
			return nil, err
		}
		return &GeoIP{db: db, lang: lang}, nil
	}
}

// Find looks up query and returns its country and area.
func (g GeoIP) Find(query string) (result fmt.Stringer, err error) {
	ip := net.ParseIP(query)
	if ip == nil {
		return nil, errors.New("query should be a valid IP")
	}
	record, err := g.db.City(ip)
	if err != nil {
		return
	}

	result = Result{
		Country:     getMapLang(record.Country.Names, g.lang),
		CountryCode: record.Country.IsoCode,
		Area:        getMapLang(record.City.Names, g.lang),
	}
	return
}

// Name returns the database name.
func (db GeoIP) Name() string {
	return "geoip"
}

// Result is a GeoIP lookup result.
type Result struct {
	Country     string `json:"country"`
	CountryCode string `json:"country_code"`
	Area        string `json:"area"`
}

func (r Result) String() string {
	if r.Area == "" {
		return r.Country
	} else {
		return fmt.Sprintf("%s %s", r.Country, r.Area)
	}
}

// DefaultLang is the fallback language used when a record lacks the requested one.
const DefaultLang = "en"

func getMapLang(data map[string]string, lang string) string {
	res, found := data[lang]
	if found {
		return res
	}
	return data[DefaultLang]
}
