package agent

import (
	"github.com/kuzhukin/metrics-collector/internal/agent/config"
	"github.com/kuzhukin/metrics-collector/internal/agent/controller"
	"github.com/kuzhukin/metrics-collector/internal/agent/reporter"
	"github.com/kuzhukin/metrics-collector/internal/zlog"
)

type Agent struct {
	ctrl *controller.Controller
}

func StartNew(config config.Config) *Agent {
	reporter := reporter.New("http://"+config.Hostport, config.TokenKey)
	agent := Agent{
		ctrl: controller.New(reporter, config.PollInterval, config.ReportInterval),
	}

	go agent.ctrl.Start()

	zlog.Logger.Infof("Metrics Agent started hostport=%v, pollinterval=%v, reportinterval=%v", config.Hostport, config.PollInterval, config.ReportInterval)

	return &agent
}

func (a *Agent) Stop() {
	zlog.Logger.Infof("Metrics Agent stopped")

	a.ctrl.Stop()
}
