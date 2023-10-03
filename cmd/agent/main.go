package main

import (
	"flag"

	"github.com/kuzhukin/metrics-collector/internal/agent"
)

const (
	hostportDefault        = "localhost:8080"
	pollIntervalSecDefault = 2
	reportIntervalDefault  = 10
)

func main() {
	if err := agent.Run(makeConfig()); err != nil {
		panic(err)
	}
}

func makeConfig() agent.Config {
	conf := agent.Config{}

	hostport := flag.String("a", hostportDefault, "Set ip:port of server")
	flag.IntVar(&conf.ReportInterval, "r", reportIntervalDefault, "Interval in seconds to sending metrics snapshot to server")
	flag.IntVar(&conf.PollingInterval, "p", pollIntervalSecDefault, "Interval in seconds for polling and collecting metrics")

	flag.Parse()

	conf.Hostport = "http://" + *hostport

	return conf
}
