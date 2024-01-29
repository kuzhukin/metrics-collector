package main

import (
	"fmt"
	"os"

	"github.com/kuzhukin/metrics-collector/internal/agent"
	"github.com/kuzhukin/metrics-collector/internal/agent/config"
	"github.com/kuzhukin/metrics-collector/internal/zlog"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	config, err := config.MakeConfig()
	if err != nil {
		return fmt.Errorf("make config, err=%w", err)
	}

	sigs := make(chan os.Signal, 1)

	metricsAgent := agent.StartNew(config)

	// waits interrupting of the agent
	sig := <-sigs

	zlog.Logger.Infof("Stop metrics agent by signal=%v\n", sig)
	metricsAgent.Stop()

	return nil
}
