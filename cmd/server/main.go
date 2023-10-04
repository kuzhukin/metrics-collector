package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

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

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT)

	srvr := server.StartNew(config)

	select {
	case sig := <-sigs:
		fmt.Printf("Stop server by signal=%v\n", sig)
		if err := srvr.Stop(); err != nil {
			return fmt.Errorf("stop server err=%s", err)
		}
	case <-srvr.Wait():
		fmt.Println("Server stopped")
	}

	return nil
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
