package ipip

import (
	"fmt"
	"log"
	"os"

	"github.com/ipipdotnet/ipdb-go"
)

// Free is an IPIP.net free (ipdb) database reader.
type Free struct {
	*ipdb.City
}

// NewIPIP opens the IPIP ipdb database at filePath.
func NewIPIP(filePath string) (*Free, error) {
	_, err := os.Stat(filePath)
	if err != nil && os.IsNotExist(err) {
		log.Printf("IPIP数据库不存在，请手动下载解压后保存到本地: %s \n", filePath)
		log.Println("下载链接： https://www.ipip.net/product/ip.html")
		return nil, err
	} else {
		db, err := ipdb.NewCity(filePath)
		if err != nil {
			return nil, err
		}
		return &Free{City: db}, nil
	}
}

// Result is an IPIP lookup result.
type Result struct {
	Country string `json:"country"`
	Region  string `json:"region"`
	City    string `json:"city"`
}

func (r Result) String() string {
	if r.City == "" {
		return fmt.Sprintf("%s %s", r.Country, r.Region)
	}
	return fmt.Sprintf("%s %s %s", r.Country, r.Region, r.City)
}

// Find looks up query and returns its country, region and city.
func (db Free) Find(query string) (result fmt.Stringer, err error) {
	info, err := db.FindInfo(query, "CN")
	if err != nil || info == nil {
		return nil, err
	} else {
		// info contains more info
		result = Result{
			Country: info.CountryName,
			Region:  info.RegionName,
			City:    info.CityName,
		}
		return
	}
}

// Name returns the database name.
func (db Free) Name() string {
	return "ipip"
}
