package config

import (
	"flag"
	"fmt"

	"github.com/caarlos0/env/v6"
)

const (
	hostportDefault        = "localhost:8080"
	storeIntervalDefault   = 300
	fileStoragePathDefault = "/tmp/metrics-db.json"
	restoreDefault         = true
)

type Config struct {
	Hostport string `env:"ADDRESS"`
	Storage  StorageConfig
}

type StorageConfig struct {
	Interval int    `env:"STORE_INTERVAL"`
	FilePath string `env:"FILE_STORAGE_PATH"`
	Restore  bool   `env:"RESTORE"`
}

func MakeConfig() (Config, error) {
	config := Config{}

	flag.StringVar(&config.Hostport, "a", hostportDefault, "Set ip:port for server")
	flag.IntVar(&config.Storage.Interval, "i", storeIntervalDefault, "Interval in secs for storing metrics to persistent storage")
	flag.StringVar(&config.Storage.FilePath, "f", fileStoragePathDefault, "Path of persistent storage")
	flag.BoolVar(&config.Storage.Restore, "r", restoreDefault, "Enable downloading metrics from persistent storage on the start")

	flag.Parse()

	if err := env.Parse(&config); err != nil {
		return config, fmt.Errorf("parse env err=%w", err)
	}

	return config, nil
}
