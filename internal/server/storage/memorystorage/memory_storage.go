package memorystorage

import (
	"fmt"

	"github.com/kuzhukin/metrics-collector/internal/metric"
	"github.com/kuzhukin/metrics-collector/internal/server/storage"
)

var _ storage.Storage = &MemoryStorage{}

type MemoryStorage struct {
	GaugeMetrics   *SyncStorage[float64]
	CounterMetrics *SyncStorage[int64]
}

func New() *MemoryStorage {
	return &MemoryStorage{
		GaugeMetrics:   NewSyncStorage[float64](),
		CounterMetrics: NewSyncStorage[int64](),
	}
}

func (s *MemoryStorage) Update(m *metric.Metric) error {
	switch m.Type {
	case metric.Gauge:
		s.GaugeMetrics.Write(m.ID, *m.Value)

		return nil
	case metric.Counter:
		s.CounterMetrics.Sum(m.ID, *m.Delta)

		return nil
	default:
		return storage.ErrUnknownKind
	}
}

func (s *MemoryStorage) BatchUpdate(metrics []*metric.Metric) error {
	for _, m := range metrics {
		if err := s.Update(m); err != nil {
			return err
		}
	}

	return nil
}

func (s *MemoryStorage) Get(kind metric.Kind, name string) (*metric.Metric, error) {
	switch kind {
	case metric.Gauge:
		gauge, ok := s.GaugeMetrics.Get(name)
		if !ok {
			return nil, fmt.Errorf("name=%s, err=%w", name, storage.ErrUnknownMetric)
		}

		return &metric.Metric{ID: name, Type: kind, Value: &gauge}, nil

	case metric.Counter:
		counter, ok := s.CounterMetrics.Get(name)
		if !ok {
			return nil, fmt.Errorf("name=%s, err=%w", name, storage.ErrUnknownMetric)
		}

		return &metric.Metric{ID: name, Type: kind, Delta: &counter}, nil
	default:
		return nil, storage.ErrUnknownKind
	}
}

func (s *MemoryStorage) List() ([]*metric.Metric, error) {
	allGauges := s.GaugeMetrics.GetAll()
	allCounters := s.CounterMetrics.GetAll()

	list := make([]*metric.Metric, 0, len(allCounters)+len(allGauges))
	list = addMetricsToList(allGauges, metric.Gauge, list)
	list = addMetricsToList(allCounters, metric.Counter, list)

	return list, nil
}

func (s *MemoryStorage) Stop() error {
	return nil
}

func addMetricsToList[T float64 | int64](metrics map[string]T, kind metric.Kind, list []*metric.Metric) []*metric.Metric {
	for name, val := range metrics {
		m, err := metric.New(name, kind, val)
		if err != nil {
			panic(err)
		}

		list = append(list, m)
	}

	return list
}
