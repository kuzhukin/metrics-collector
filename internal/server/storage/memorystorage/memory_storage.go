package memorystorage

import (
	"errors"

	"github.com/kuzhukin/metrics-collector/internal/server/metric"
	"github.com/kuzhukin/metrics-collector/internal/server/storage"
)

var ErrUnknownMetric = errors.New("unknown metric name")
var ErrUnknownKind = errors.New("unknown metric kind")

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
		s.gaugeMetrics.Write(m.Name, m.Value.Gauge())

		return nil
	case metric.Counter:
		s.counterMetrics.Sum(m.Name, m.Value.Counter())

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

func (s *memoryStorage) List() []*metric.Metric {
	allGauges := s.gaugeMetrics.GetAll()
	allCounters := s.counterMetrics.GetAll()

	list := make([]*metric.Metric, 0, len(allCounters)+len(allGauges))
	list = addMetricsToListp(allGauges, metric.Gauge, list)
	list = addMetricsToListp(allCounters, metric.Counter, list)

	return list
}

func addMetricsToListp[T float64 | int64](metrics map[string]T, kind metric.Kind, list []*metric.Metric) []*metric.Metric {
	for name, valT := range metrics {
		list = append(list, metric.NewMetric(kind, name, metric.NewValueByKind(kind, valT)))
	}

	return list
}
