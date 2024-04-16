package main

import (
	"fmt"
	"os"

	"go.uber.org/zap"

	"github.com/nkitlabs/go-http-gorm-example/pkg/config"
	dbstore "github.com/nkitlabs/go-http-gorm-example/pkg/db"
)

func main() {
	logger := zap.Must(zap.NewProduction())
	defer func() {
		if err := logger.Sync(); err != nil {
			panic(err)
		}
	}()

	config_filename := os.Getenv("CONFIG_FILE")
	if config_filename == "" {
		config_filename = "config.dev.yaml"
	}

	logger.Info(fmt.Sprintf("Read config file from %s", config_filename))
	conf, err := config.ReadConfig(config_filename)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	db, err := dbstore.Init(conf.Conn)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	_ = db
}
