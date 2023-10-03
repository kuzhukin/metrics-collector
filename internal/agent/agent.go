package agent

import (
	"github.com/kuzhukin/metrics-collector/internal/agent/controller"
	"github.com/kuzhukin/metrics-collector/internal/agent/reporter"
)

const hostport = "http://localhost:8080"

func Run() error {
	reporter := reporter.New(hostport)
	controller.New(reporter).Start()

	return nil
}
