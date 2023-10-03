package main

import "github.com/kuzhukin/metrics-collector/internal/agent"

const hostport = "http://localhost:8080"

func main() {
	if err := agent.Run(); err != nil {
		panic(err)
	}
}
