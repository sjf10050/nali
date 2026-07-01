package config

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
	"github.com/zu1k/nali/internal/db"
)

func ReadConfig(basePath string) error {
	viper.SetDefault("databases", db.GetDefaultDBList())
	viper.SetDefault("selected.ipv4", "qqwry")
	viper.SetDefault("selected.ipv6", "zxipv6wry")
	viper.SetDefault("selected.cdn", "cdn")
	viper.SetDefault("selected.lang", "zh-CN")

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(basePath)
	err := viper.ReadInConfig()
	if err != nil {
		err = viper.SafeWriteConfig()
		if err != nil {
			return fmt.Errorf("failed to write config: %w", err)
		}
	}

	if err := viper.BindEnv("selected.ipv4", "NALI_DB_IP4"); err != nil {
		log.Println("bind env selected.ipv4:", err)
	}
	if err := viper.BindEnv("selected.ipv6", "NALI_DB_IP6"); err != nil {
		log.Println("bind env selected.ipv6:", err)
	}
	if err := viper.BindEnv("selected.cdn", "NALI_DB_CDN"); err != nil {
		log.Println("bind env selected.cdn:", err)
	}
	if err := viper.BindEnv("selected.lang", "NALI_LANG"); err != nil {
		log.Println("bind env selected.lang:", err)
	}

	dbList := db.List{}
	err = viper.UnmarshalKey("databases", &dbList)
	if err != nil {
		return fmt.Errorf("config invalid: %w", err)
	}

	db.NameDBMap.From(dbList)
	db.TypeDBMap.From(dbList)

	return nil
}
