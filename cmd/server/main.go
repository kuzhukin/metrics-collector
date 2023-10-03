package main

import (
	"github.com/kuzhukin/metrics-collector/internal/server"
)

func main() {
	if err := server.Run(); err != nil {
		panic(err)
	}
}
