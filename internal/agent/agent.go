package agent

import (
	"fmt"

	"github.com/kuzhukin/metrics-collector/internal/agent/config"
	"github.com/kuzhukin/metrics-collector/internal/agent/controller"
	"github.com/kuzhukin/metrics-collector/internal/agent/reporter"
	"github.com/kuzhukin/metrics-collector/internal/zlog"
)

type Agent struct {
	ctrl *controller.Controller
}

// StartNew - creats and starts new metrics agent
func StartNew(config config.Config) (*Agent, error) {
	reporter, err := reporter.New(config)
	if err != nil {
		return nil, fmt.Errorf("new reporter, err=%w", err)
	}

	agent := Agent{
		ctrl: controller.New(reporter, config.PollInterval, config.ReportInterval),
	}

	go agent.ctrl.Start()

	zlog.Logger.Infof("Metrics Agent started  config=%+v", config)

	return &agent, nil
}

func (a *Agent) Stop() {
	zlog.Logger.Infof("Metrics Agent stopped")

	a.ctrl.Stop()
}
