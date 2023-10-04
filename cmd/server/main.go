package main

import (
	"flag"
	"fmt"

	"github.com/caarlos0/env/v6"
	"github.com/kuzhukin/metrics-collector/internal/server"
)

const hostportDefault = "localhost:8080"

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

	return server.Run(config)
}
func makeConfig() (server.Config, error) {
	config := server.Config{}

	flag.StringVar(&config.Hostport, "a", hostportDefault, "Set ip:port for server")
	flag.Parse()

	if err := env.Parse(&config); err != nil {
		return config, fmt.Errorf("parse env err=%w", err)
	}

	return config, nil
}
