package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/kuzhukin/metrics-collector/internal/server"
	"github.com/kuzhukin/metrics-collector/internal/zlog"
)

var (
	buildVersion = "N/A"
	buildDate    = "N/A"
	buildCommit  = "N/A"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	defer func() {
		// flush logs
		_ = zlog.Logger.Sync()
	}()

	printBuildInfo()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT)

	// starting new metrics HTTP server
	srvr, err := server.StartNew()
	if err != nil {
		return fmt.Errorf("server start err=%w", err)
	}

	// waits interrupting or stop of the server
	select {
	case sig := <-sigs:
		zlog.Logger.Infof("Stop server by signal=%v\n", sig)
		if err := srvr.Stop(); err != nil {
			return fmt.Errorf("stop server err=%s", err)
		}
	case <-srvr.Wait():
		zlog.Logger.Info("Server stopped")
	}

	return nil
}

func printBuildInfo() {
	if buildVersion == "" {
		buildVersion = "N/A"
	}

	if buildDate == "" {
		buildDate = "N/A"
	}

	if buildCommit == "" {
		buildCommit = "N/A"
	}

	fmt.Printf(
		`Build version: %s (или "N/A" при отсутствии значения)\n`+
			`Build date: %s (или "N/A" при отсутствии значения)\n`+
			`Build commit: %s (или "N/A" при отсутствии значения)\n`,
		buildVersion, buildDate, buildCommit,
	)
}
