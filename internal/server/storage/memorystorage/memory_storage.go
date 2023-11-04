package memorystorage

import (
	"errors"
	"fmt"

	"github.com/kuzhukin/metrics-collector/internal/server/metric"
	"github.com/kuzhukin/metrics-collector/internal/server/storage"
	"github.com/kuzhukin/metrics-collector/internal/zlog"
)

var ErrUnknownMetric = errors.New("unknown metric name")
var ErrUnknownKind = errors.New("unknown metric kind")

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
	switch m.Kind {
	case metric.Gauge:
		s.GaugeMetrics.Write(m.Name, m.Value.Gauge())

		return nil
	case metric.Counter:
		s.CounterMetrics.Sum(m.Name, m.Value.Counter())

		return nil
	default:
		return ErrUnknownKind
	}
}

func (s *MemoryStorage) Get(kind metric.Kind, name string) (*metric.Metric, error) {
	switch kind {
	case metric.Gauge:
		gauge, ok := s.GaugeMetrics.Get(name)
		if !ok {
			return nil, fmt.Errorf("name=%s, err=%w", name, ErrUnknownMetric)
		}

		return metric.NewMetric(kind, name, metric.GaugeValue(gauge)), nil

	case metric.Counter:
		counter, ok := s.CounterMetrics.Get(name)
		if !ok {
			return nil, fmt.Errorf("name=%s, err=%w", name, ErrUnknownMetric)
		}

		return metric.NewMetric(kind, name, metric.CounterValue(counter)), nil
	default:
		return nil, ErrUnknownKind
	}
}

func (s *MemoryStorage) List() []*metric.Metric {
	allGauges := s.GaugeMetrics.GetAll()
	allCounters := s.CounterMetrics.GetAll()

	list := make([]*metric.Metric, 0, len(allCounters)+len(allGauges))
	list = addMetricsToList(allGauges, metric.Gauge, list)
	list = addMetricsToList(allCounters, metric.Counter, list)

	return list
}

func (s *MemoryStorage) Stop() error {
	return nil
}

func addMetricsToList[T float64 | int64](metrics map[string]T, kind metric.Kind, list []*metric.Metric) []*metric.Metric {
	for name, valT := range metrics {
		val, err := metric.NewValueByKind(kind, valT)
		if err != nil {
			zlog.Logger.Errorf("new value by kind=%s, err=%s", kind, err)
			continue
		}

		list = append(list, metric.NewMetric(kind, name, val))
	}

	return list
}
