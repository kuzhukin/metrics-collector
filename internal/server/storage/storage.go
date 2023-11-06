package storage

import (
	"errors"

	"github.com/kuzhukin/metrics-collector/internal/server/metric"
)

var ErrUnknownMetric = errors.New("unknown metric name")
var ErrUnknownKind = errors.New("unknown metric kind")

//go:generate mockery --name=Storage --filename=storage.go --outpkg=mockstorage --output=mockstorage
type Storage interface {
	Update(m *metric.Metric) error
	Get(kind metric.Kind, name string) (*metric.Metric, error)
	List() ([]*metric.Metric, error)
}
