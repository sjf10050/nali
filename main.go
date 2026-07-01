package main

import (
	"log"

	"github.com/zu1k/nali/cmd"
	"github.com/zu1k/nali/internal/config"
	"github.com/zu1k/nali/internal/constant"
	"github.com/zu1k/nali/internal/db"
	"github.com/zu1k/nali/internal/migration"
)

func main() {
	if err := constant.InitPaths(); err != nil {
		log.Fatalln("Failed to prepare config/data directories:", err)
	}
	migration.Run()
	if err := config.ReadConfig(constant.ConfigDirPath); err != nil {
		log.Fatalln("Failed to read config:", err)
	}
	db.PreloadConfig()
	cmd.Execute()
}
