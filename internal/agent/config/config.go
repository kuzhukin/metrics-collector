package config

import (
	"flag"
	"fmt"

	"github.com/caarlos0/env/v6"
)

const (
	hostportDefault        = "localhost:8080"
	pollIntervalSecDefault = 2
	reportIntervalDefault  = 10
)

type Config struct {
	Hostport       string `env:"ADDRESS"`
	ReportInterval int    `env:"REPORT_INTERVAL"`
	PollInterval   int    `env:"POLL_INTERVAL"`
	SingnatureKey  string `env:"KEY"`
}

func MakeConfig() (Config, error) {
	config := Config{}

	flag.StringVar(&config.Hostport, "a", hostportDefault, "Set ip:port of server")
	flag.IntVar(&config.ReportInterval, "r", reportIntervalDefault, "Interval in seconds for sending metrics snapshot to server")
	flag.IntVar(&config.PollInterval, "p", pollIntervalSecDefault, "Interval in seconds for polling and collecting metrics")
	flag.StringVar(&config.SingnatureKey, "k", "", "Signature key")
	flag.Parse()

	if err := env.Parse(&config); err != nil {
		return config, fmt.Errorf("parse env err=%w", err)
	}

	return config, nil
}
