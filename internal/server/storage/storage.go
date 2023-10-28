package storage

import "github.com/kuzhukin/metrics-collector/internal/server/metric"

//go:generate mockery --name=Storage --filename=storage.go --outpkg=mockstorage --output=mockstorage
type Storage interface {
	// Initialize() error
	Update(m *metric.Metric) error
	Get(kind metric.Kind, name string) (*metric.Metric, error)
	List() []*metric.Metric
}
