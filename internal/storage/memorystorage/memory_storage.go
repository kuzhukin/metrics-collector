package memorystorage

import (
	"fmt"

	"github.com/kuzhukin/metrics-collector/internal/metric"
	"github.com/kuzhukin/metrics-collector/internal/storage"
)

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
		return s.updateGauge(m.Name, m.Value.Gauge())
	case metric.Counter:
		return s.updateCounter(m.Name, m.Value.Counter())
	default:
		return fmt.Errorf("doesn't have update handle func for kind=%s", m.Kind)
	}
}

func (s *memoryStorage) updateGauge(name string, value float64) error {
	s.gaugeMetrics.Update(func(storage map[string]float64) {
		storage[name] = value
	})

	return nil
}

func (s *memoryStorage) updateCounter(name string, value int64) error {
	s.counterMetrics.Update(func(storage map[string]int64) {
		storage[name] = value
	})

	return nil
}
