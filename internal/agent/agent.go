package agent

import (
	"github.com/kuzhukin/metrics-collector/internal/agent/controller"
	"github.com/kuzhukin/metrics-collector/internal/agent/reporter"
)

type Agent struct {
	ctrl *controller.Controller
}

func StartNew(config Config) *Agent {
	reporter := reporter.New("http://" + config.Hostport)
	agent := Agent{ctrl: controller.New(reporter, config.PollInterval, config.ReportInterval)}

	go agent.ctrl.Start()

	return &agent
}

func (a *Agent) Stop() {
	a.ctrl.Stop()
}
