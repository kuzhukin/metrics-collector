package memorystorage

import (
	"errors"

	"github.com/kuzhukin/metrics-collector/internal/metric"
	"github.com/kuzhukin/metrics-collector/internal/storage"
)

var ErrUnknownMetric = errors.New("Unknown metric name")
var ErrUnknownKind = errors.New("Unknown metric kind")

var _ storage.Storage = &memoryStorage{}

type memoryStorage struct {
	gaugeMetrics   syncMemoryStorage[float64]
	counterMetrics syncMemoryStorage[int64]
}

func New() storage.Storage {
	return &memoryStorage{
		gaugeMetrics:   newSyncMemoryStorage[float64](),
		counterMetrics: newSyncMemoryStorage[int64](),
	}
}

func (s *memoryStorage) Update(m *metric.Metric) error {
	switch m.Kind {
	case metric.Gauge:
		s.gaugeMetrics.Update(m.Name, m.Value.Gauge())

		return nil
	case metric.Counter:
		s.counterMetrics.Update(m.Name, m.Value.Counter())

		return nil
	default:
		return ErrUnknownKind
	}
}

func (s *memoryStorage) Get(kind metric.Kind, name string) (*metric.Metric, error) {
	switch kind {
	case metric.Gauge:
		gauge, ok := s.gaugeMetrics.Get(name)
		if !ok {
			return nil, ErrUnknownMetric
		}

		return metric.NewMetric(kind, name, metric.GaugeValue(gauge)), nil

	case metric.Counter:
		counter, ok := s.counterMetrics.Get(name)
		if !ok {
			return nil, ErrUnknownMetric
		}

		return metric.NewMetric(kind, name, metric.CounterValue(counter)), nil
	default:
		return nil, ErrUnknownKind
	}
}
