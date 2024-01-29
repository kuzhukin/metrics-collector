// module declares the metrics storage interface
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
	// creates or updates a metric in the storage
	Update(ctx context.Context, m *metric.Metric) error
	// creates or updates a batch of metrics in the storage
	BatchUpdate(ctx context.Context, m []*metric.Metric) error
	// returns a metric by kind and name
	Get(ctx context.Context, kind metric.Kind, name string) (*metric.Metric, error)
	// returns all metrics from the storage
	List(ctx context.Context) ([]*metric.Metric, error)
}
