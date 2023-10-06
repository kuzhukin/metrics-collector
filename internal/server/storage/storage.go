package storage

import "github.com/kuzhukin/metrics-collector/internal/server/metric"

//go:generate mockery --name=Storage --outpkg=mockstorage --output=mockstorage
type Storage interface {
	Update(m *metric.Metric) error
	Get(kind metric.Kind, name string) (*metric.Metric, error)
	List() []*metric.Metric
}
