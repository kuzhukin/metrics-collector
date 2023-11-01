package agent

import (
	"github.com/kuzhukin/metrics-collector/internal/agent/controller"
	"github.com/kuzhukin/metrics-collector/internal/agent/reporter"
	"github.com/kuzhukin/metrics-collector/internal/log"
)

type Agent struct {
	ctrl *controller.Controller
}

func StartNew(config Config) *Agent {
	reporter := reporter.New("http://" + config.Hostport)
	agent := Agent{
		ctrl: controller.New(reporter, config.PollInterval, config.ReportInterval),
	}

	go agent.ctrl.Start()

	log.Logger.Infof("Metrics Agent started hostport=%v, pollinterval=%v, reportinterval=%v", config.Hostport, config.PollInterval, config.ReportInterval)

	return &agent
}

func (a *Agent) Stop() {
	log.Logger.Infof("Metrics Agent stopped")

	a.ctrl.Stop()
}
