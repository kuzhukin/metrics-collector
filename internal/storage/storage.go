package storage

import "github.com/kuzhukin/metrics-collector/internal/metric"

//go:generate mockgen -source=storage.go -destination=mockstorage/mock.go -package=mockstorage
type Storage interface {
	Update(m *metric.Metric) error
	Get(kind metric.Kind, name string) (*metric.Metric, error)
}
