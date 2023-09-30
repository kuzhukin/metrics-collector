package main

import (
	"github.com/kuzhukin/metrics-collector/cmd/agent/controller"
	"github.com/kuzhukin/metrics-collector/cmd/agent/reporter"
)

const hostport = "http://localhost:8080"

func main() {
	reporter := reporter.New(hostport)
	controller.New(reporter).Start()
}
