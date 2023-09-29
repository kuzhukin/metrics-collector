package memorystorage

import "github.com/kuzhukin/metrics-collector/internal/storage"

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

func (s *memoryStorage) UpdateGauge(name string, value float64) error {
	s.gaugeMetrics.Update(func(storage map[string]float64) {
		storage[name] = value
	})

	return nil
}

func (s *memoryStorage) UpdateCounter(name string, value int64) error {
	s.counterMetrics.Update(func(storage map[string]int64) {
		storage[name] = value
	})

	return nil
}
