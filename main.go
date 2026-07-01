package main

import (
	"log"

	"github.com/zu1k/nali/internal/constant"

	"github.com/zu1k/nali/cmd"
	"github.com/zu1k/nali/internal/config"

	_ "github.com/zu1k/nali/internal/migration"
)

func main() {
	if err := config.ReadConfig(constant.ConfigDirPath); err != nil {
		log.Fatalln("Failed to read config:", err)
	}
	cmd.Execute()
}
