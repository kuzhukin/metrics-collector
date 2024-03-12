package reporter

import "github.com/kuzhukin/metrics-collector/internal/agent/config"

const batchUpdateEndpoint = "/updates/"

// Reporter sends metrics to server
//
//go:generate mockery --name=Reporter --filename=reporter.go --outpkg=mockreporter --output=mockreporter
type Reporter interface {
	Report(gaugeMetrics map[string]float64, counterMetrics map[string]int64)
}

func New(config config.Config) (Reporter, error) {
	if config.UseGRPC {
		return newGRPCReporter(config)
	}

	return newHTTPReporter(config)
}
