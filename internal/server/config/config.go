package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/caarlos0/env/v6"
	"github.com/kuzhukin/metrics-collector/internal/zlog"
)

const (
	hostportDefault        = "localhost:8080"
	storeIntervalDefault   = 300
	fileStoragePathDefault = "/tmp/metrics-db.json"
	restoreDefault         = true
)

// Config of HTTP server
type Config struct {
	// listening address:port
	Hostport string `env:"ADDRESS" json:"address"`
	// key for signature calculation
	SingnatureKey string `env:"KEY" json:"key"`
	// storage config
	Storage StorageConfig
	// flag for enabling logger
	EnableLogger bool `env:"ENABLE_LOGGER" json:"enable_logger"`
	// key for https connection
	CryptoKey string `env:"CRYPTO_KEY" json:"crypto_key"`
	// path to config file
	ConfigFilePath string `env:"CONFIG"`
}

// StorageConfig - metrics storage config
type StorageConfig struct {
	// path to the file (for file storage)
	FilePath string `env:"FILE_STORAGE_PATH" json:"store_file"`
	// dsn for database storage
	DatabaseDSN string `env:"DATABASE_DSN" json:"database_dsn"`
	// interval of updloading metrics to persistent storage
	Interval int `env:"STORE_INTERVAL" json:"store_interval"`
	// enable downloading metrics from persistent storage on server start
	Restore bool `env:"RESTORE" json:"restore"`
}

// MakeConfig - reads configuration from application parameters and environment variables
func MakeConfig() (Config, error) {
	config := Config{}

	flag.StringVar(&config.Hostport, "a", hostportDefault, "Set ip:port for server")
	flag.IntVar(&config.Storage.Interval, "i", storeIntervalDefault, "Interval in secs for storing metrics to persistent storage")
	flag.StringVar(&config.Storage.FilePath, "f", fileStoragePathDefault, "Path of persistent storage")
	flag.BoolVar(&config.Storage.Restore, "r", restoreDefault, "Enable downloading metrics from persistent storage on the start")
	flag.StringVar(&config.Storage.DatabaseDSN, "d", "", "Database connection string")
	flag.StringVar(&config.SingnatureKey, "k", "", "Signature key")
	flag.StringVar(&config.CryptoKey, "crypto-key", "", "Crypto key")
	flag.BoolVar(&config.EnableLogger, "l", true, "Enable logger")

	flag.Parse()

	if err := env.Parse(&config); err != nil {
		return config, fmt.Errorf("parse env err=%w", err)
	}

	return updateConfigFromFile(config), nil
}

func updateConfigFromFile(config Config) Config {
	if config.ConfigFilePath != "" {
		jsonConfig := Config{}

		data, err := os.ReadFile(config.ConfigFilePath)
		if err != nil {
			zlog.Logger.Errorf("error read config from file=%s, err=%s", config.ConfigFilePath, err)
		} else {
			if err = json.Unmarshal(data, &jsonConfig); err != nil {
				zlog.Logger.Errorf("unmarshal config from file=%s, err=%s", config.ConfigFilePath, err)
			} else {
				if config.Hostport == "" {
					config.Hostport = jsonConfig.Hostport
				}

				if config.SingnatureKey == "" {
					config.SingnatureKey = jsonConfig.SingnatureKey
				}

				if config.Storage.FilePath == "" {
					config.Storage.FilePath = jsonConfig.Storage.FilePath
				}

				if config.Storage.DatabaseDSN == "" {
					config.Storage.DatabaseDSN = jsonConfig.Storage.DatabaseDSN
				}

				if config.Storage.Interval == 0 {
					config.Storage.Interval = jsonConfig.Storage.Interval
				}

				if !config.Storage.Restore {
					config.Storage.Restore = jsonConfig.Storage.Restore
				}

				if !config.EnableLogger {
					config.EnableLogger = jsonConfig.EnableLogger
				}

				if config.CryptoKey == "" {
					config.CryptoKey = jsonConfig.CryptoKey
				}
			}
		}
	}

	return config
}
