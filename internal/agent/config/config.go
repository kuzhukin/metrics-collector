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
	// server address:port for reporting metrics
	Hostport string `env:"ADDRESS"`
	// key for payload signature
	SingnatureKey string `env:"KEY"`
	// interval of metrics reporting
	ReportInterval int `env:"REPORT_INTERVAL"`
	// interval of polling and collecting metrics
	PollInterval int `env:"POLL_INTERVAL"`
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
