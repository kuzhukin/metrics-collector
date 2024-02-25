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

// StartNew - creats and starts new metrics agent
func StartNew(config config.Config) *Agent {
	reporter := reporter.New("http://"+config.Hostport, config.SingnatureKey)
	agent := Agent{
		ctrl: controller.New(reporter, config.PollInterval, config.ReportInterval),
	}

	go agent.ctrl.Start()

	zlog.Logger.Infof("Metrics Agent started  config=%+v", config)

	return &agent
}

func (a *Agent) Stop() {
	zlog.Logger.Infof("Metrics Agent stopped")

	a.ctrl.Stop()
}
