package storage

import (
	"context"
	"errors"

	"github.com/kuzhukin/metrics-collector/internal/metric"
)

var ErrUnknownMetric = errors.New("unknown metric name")
var ErrUnknownKind = errors.New("unknown metric kind")

//go:generate mockery --name=Storage --filename=storage.go --outpkg=mockstorage --output=mockstorage
type Storage interface {
	Update(ctx context.Context, m *metric.Metric) error
	BatchUpdate(ctx context.Context, m []*metric.Metric) error
	Get(ctx context.Context, kind metric.Kind, name string) (*metric.Metric, error)
	List(ctx context.Context) ([]*metric.Metric, error)
}
