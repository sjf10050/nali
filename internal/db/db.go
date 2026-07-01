package db

import (
	"fmt"
	"log"
	"net"

	"github.com/spf13/viper"

	"github.com/zu1k/nali/pkg/cdn"
	"github.com/zu1k/nali/pkg/dbif"
	"github.com/zu1k/nali/pkg/geoip"
	"github.com/zu1k/nali/pkg/qqwry"
	"github.com/zu1k/nali/pkg/zxipv6wry"
)

var nat64CIDR *net.IPNet

func init() {
	_, nat64CIDR, _ = net.ParseCIDR("64:ff9b::/96")
}

func GetDB(typ dbif.QueryType) (db dbif.DB, err error) {
	cacheMu.RLock()
	cached, found := dbTypeCache[typ]
	cacheMu.RUnlock()
	if found {
		return cached, nil
	}

	lang := viper.GetString("selected.lang")
	if lang == "" {
		lang = "zh-CN"
	}

	switch typ {
	case dbif.TypeIPv4:
		selected := viper.GetString("selected.ipv4")
		if selected != "" {
			dbInfo, dbErr := getDbByName(selected)
			if dbErr != nil {
				return nil, dbErr
			}
			db, err = dbInfo.get()
			break
		}

		if lang == "zh-CN" {
			qqwryDB, dbErr := getDbByName("qqwry")
			if dbErr != nil {
				return nil, dbErr
			}
			db, err = qqwry.NewQQwry(qqwryDB.File)
		} else {
			geoipDB, dbErr := getDbByName("geoip")
			if dbErr != nil {
				return nil, dbErr
			}
			db, err = geoip.NewGeoIP(geoipDB.File)
		}
	case dbif.TypeIPv6:
		selected := viper.GetString("selected.ipv6")
		if selected != "" {
			dbInfo, dbErr := getDbByName(selected)
			if dbErr != nil {
				return nil, dbErr
			}
			db, err = dbInfo.get()
			break
		}

		if lang == "zh-CN" {
			zxDB, dbErr := getDbByName("zxipv6wry")
			if dbErr != nil {
				return nil, dbErr
			}
			db, err = zxipv6wry.NewZXwry(zxDB.File)
		} else {
			geoipDB, dbErr := getDbByName("geoip")
			if dbErr != nil {
				return nil, dbErr
			}
			db, err = geoip.NewGeoIP(geoipDB.File)
		}
	case dbif.TypeDomain:
		selected := viper.GetString("selected.cdn")
		if selected != "" {
			dbInfo, dbErr := getDbByName(selected)
			if dbErr != nil {
				return nil, dbErr
			}
			db, err = dbInfo.get()
			break
		}

		cdnDB, dbErr := getDbByName("cdn")
		if dbErr != nil {
			return nil, dbErr
		}
		db, err = cdn.NewCDN(cdnDB.File)
	default:
		return nil, fmt.Errorf("query type not supported: %v", typ)
	}

	if err != nil {
		return nil, err
	}
	if db == nil {
		return nil, fmt.Errorf("database init failed")
	}

	cacheMu.Lock()
	dbTypeCache[typ] = db
	cacheMu.Unlock()
	return db, nil
}

func Find(typ dbif.QueryType, query string) *Result {
	if cached, found := queryCache.Load(query); found {
		result, ok := cached.(*Result)
		if ok {
			return result
		}
	}
	// Convert NAT64 64:ff9b::/96 to IPv4
	if typ == dbif.TypeIPv6 {
		ip := net.ParseIP(query)
		if ip != nil && nat64CIDR != nil && nat64CIDR.Contains(ip) {
			ip4 := make(net.IP, 4)
			copy(ip4, ip[12:16])
			query = ip4.String()
			typ = dbif.TypeIPv4
		}
	}
	db, err := GetDB(typ)
	if err != nil {
		log.Println("GetDB error:", err)
		return nil
	}
	result, err := db.Find(query)
	if err != nil {
		return nil
	}
	res := &Result{db.Name(), result}
	queryCache.Store(query, res)
	return res
}
