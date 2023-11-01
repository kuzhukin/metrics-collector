package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/kuzhukin/metrics-collector/internal/log"
	"github.com/kuzhukin/metrics-collector/internal/server"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	defer func() {
		_ = log.Logger.Sync()
	}()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT)

	srvr, err := server.StartNew()
	if err != nil {
		return fmt.Errorf("server start err=%w", err)
	}

	select {
	case sig := <-sigs:
		log.Logger.Infof("Stop server by signal=%v\n", sig)
		if err := srvr.Stop(); err != nil {
			return fmt.Errorf("stop server err=%s", err)
		}
	case <-srvr.Wait():
		log.Logger.Info("Server stopped")
	}

	return nil
}
