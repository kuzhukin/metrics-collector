package main

import (
	"flag"

	"github.com/kuzhukin/metrics-collector/internal/server"
)

const hostportDefault = "localhost:8080"

func main() {
	if err := server.Run(makeConfig()); err != nil {
		panic(err)
	}
}

func makeConfig() server.Config {
	conf := server.Config{}

	flag.StringVar(&conf.Hostport, "a", hostportDefault, "Set ip:port for server")
	flag.Parse()

	return conf
}
