package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/caarlos0/env/v6"
	"github.com/kuzhukin/metrics-collector/internal/agent"
	"github.com/kuzhukin/metrics-collector/internal/log"
)

const (
	hostportDefault        = "localhost:8080"
	pollIntervalSecDefault = 2
	reportIntervalDefault  = 10
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	config, err := makeConfig()
	if err != nil {
		return fmt.Errorf("make config, err=%w", err)
	}

	sigs := make(chan os.Signal, 1)

	metricsAgent := agent.StartNew(config)

	sig := <-sigs

	log.Logger.Infof("Stop metrics agent by signal=%v\n", sig)
	metricsAgent.Stop()

	return nil
}

func makeConfig() (agent.Config, error) {
	config := agent.Config{}

	flag.StringVar(&config.Hostport, "a", hostportDefault, "Set ip:port of server")
	flag.IntVar(&config.ReportInterval, "r", reportIntervalDefault, "Interval in seconds for sending metrics snapshot to server")
	flag.IntVar(&config.PollInterval, "p", pollIntervalSecDefault, "Interval in seconds for polling and collecting metrics")
	flag.Parse()

	if err := env.Parse(&config); err != nil {
		return config, fmt.Errorf("parse env err=%w", err)
	}

	return config, nil
}
