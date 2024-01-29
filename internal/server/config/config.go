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

// HTTP server config
type Config struct {
	// listening address:port
	Hostport string `env:"ADDRESS"`
	// key for signature calculation
	SingnatureKey string `env:"KEY"`
	// flag for enabling logger
	EnableLogger bool `env:"ENABLE_LOGGER"`
	// storage config
	Storage StorageConfig
}

type StorageConfig struct {
	// interval of updloading metrics to persistent storage
	Interval int `env:"STORE_INTERVAL"`
	// path to the file (for file storage)
	FilePath string `env:"FILE_STORAGE_PATH"`
	// enable downloading metrics from persistent storage on server start
	Restore bool `env:"RESTORE"`
	// dsn for database storage
	DatabaseDSN string `env:"DATABASE_DSN"`
}

// reading configuration from application parameters and environment variables
func MakeConfig() (Config, error) {
	config := Config{}

	flag.StringVar(&config.Hostport, "a", hostportDefault, "Set ip:port for server")
	flag.IntVar(&config.Storage.Interval, "i", storeIntervalDefault, "Interval in secs for storing metrics to persistent storage")
	flag.StringVar(&config.Storage.FilePath, "f", fileStoragePathDefault, "Path of persistent storage")
	flag.BoolVar(&config.Storage.Restore, "r", restoreDefault, "Enable downloading metrics from persistent storage on the start")
	flag.StringVar(&config.Storage.DatabaseDSN, "d", "", "Database connection string")
	flag.StringVar(&config.SingnatureKey, "k", "", "Signature key")
	flag.BoolVar(&config.EnableLogger, "l", true, "Enable logger")

	flag.Parse()

	if err := env.Parse(&config); err != nil {
		return config, fmt.Errorf("parse env err=%w", err)
	}

	return config, nil
}
