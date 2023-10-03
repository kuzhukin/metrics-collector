package agent

import (
	"github.com/kuzhukin/metrics-collector/internal/agent/controller"
	"github.com/kuzhukin/metrics-collector/internal/agent/reporter"
)

func Run(config Config) error {
	reporter := reporter.New(config.Hostport)
	controller.New(reporter, config.PollingInterval, config.ReportInterval).Start()

	return nil
}
