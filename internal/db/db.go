package db

import (
	"fmt"
	"log"
	"net"

	"github.com/zu1k/nali/internal/constant"
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

	switch typ {
	case dbif.TypeIPv4:
		selected := selectedIPv4
		if selected != "" {
			dbInfo, dbErr := getDbByName(selected)
			if dbErr != nil {
				return nil, dbErr
			}
			db, err = dbInfo.get()
			break
		}

		if selectedLang == "zh-CN" {
			qqwryDB, dbErr := getDbByName("qqwry")
			if dbErr != nil {
				return nil, dbErr
			}
			db, err = qqwry.NewQQwry(constant.ResolveDBPath(qqwryDB.File))
		} else {
			geoipDB, dbErr := getDbByName("geoip")
			if dbErr != nil {
				return nil, dbErr
			}
			db, err = geoip.NewGeoIP(constant.ResolveDBPath(geoipDB.File), selectedLang)
		}
	case dbif.TypeIPv6:
		selected := selectedIPv6
		if selected != "" {
			dbInfo, dbErr := getDbByName(selected)
			if dbErr != nil {
				return nil, dbErr
			}
			db, err = dbInfo.get()
			break
		}

		if selectedLang == "zh-CN" {
			zxDB, dbErr := getDbByName("zxipv6wry")
			if dbErr != nil {
				return nil, dbErr
			}
			db, err = zxipv6wry.NewZXwry(constant.ResolveDBPath(zxDB.File))
		} else {
			geoipDB, dbErr := getDbByName("geoip")
			if dbErr != nil {
				return nil, dbErr
			}
			db, err = geoip.NewGeoIP(constant.ResolveDBPath(geoipDB.File), selectedLang)
		}
	case dbif.TypeDomain:
		selected := selectedCDN
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
		db, err = cdn.NewCDN(constant.ResolveDBPath(cdnDB.File))
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
	// Cache identity is the original (typ, query). The type is part of the key so
	// the same text queried as different types can't collide, and the NAT64
	// rewrite below changes only the lookup target — not the key — so NAT64
	// inputs still cache under their own text.
	cacheKey := fmt.Sprintf("%d\x00%s", typ, query)
	if result, found := queryCache.Load(cacheKey); found {
		return result
	}

	// Convert NAT64 64:ff9b::/96 to IPv4 for the lookup only.
	lookupTyp, lookupQuery := typ, query
	if typ == dbif.TypeIPv6 {
		ip := net.ParseIP(query)
		if ip != nil && nat64CIDR != nil && nat64CIDR.Contains(ip) {
			ip4 := make(net.IP, 4)
			copy(ip4, ip[12:16])
			lookupQuery = ip4.String()
			lookupTyp = dbif.TypeIPv4
		}
	}

	db, err := GetDB(lookupTyp)
	if err != nil {
		log.Println("GetDB error:", err)
		return nil
	}
	result, err := db.Find(lookupQuery)
	if err != nil {
		return nil
	}
	res := &Result{db.Name(), result}
	queryCache.Store(cacheKey, res)
	return res
}
