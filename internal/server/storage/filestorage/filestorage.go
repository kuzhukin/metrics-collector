package filestorage

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/kuzhukin/metrics-collector/internal/server/config"
	"github.com/kuzhukin/metrics-collector/internal/server/metric"
	"github.com/kuzhukin/metrics-collector/internal/server/storage"
	"github.com/kuzhukin/metrics-collector/internal/server/storage/memorystorage"
	"github.com/kuzhukin/metrics-collector/internal/transport"
	"github.com/kuzhukin/metrics-collector/internal/zlog"
)

var _ storage.Storage = &FileStorage{}

type FileStorage struct {
	memoryStorage memorystorage.MemoryStorage

	filepath string
	interval time.Duration
}

func New(config config.StorageConfig) (*FileStorage, error) {
	storage := &FileStorage{
		memoryStorage: memorystorage.MemoryStorage{
			GaugeMetrics:   memorystorage.NewSyncStorage[float64](),
			CounterMetrics: memorystorage.NewSyncStorage[int64](),
		},

		filepath: config.FilePath,
		interval: time.Second * time.Duration(config.Interval),
	}

	if config.Restore {
		if err := storage.restore(); err != nil {
			zlog.Logger.Warnf("Restore metrics from file=%v err=%s", storage.filepath, err)
		}
	}

	if storage.interval > 0 {
		storage.startSyncer()
	}

	return storage, nil
}

func (s *FileStorage) Update(m *metric.Metric) error {
	if err := s.memoryStorage.Update(m); err != nil {
		return err
	}

	if s.interval != 0 {
		return nil
	}

	return s.sync()
}

func (s *FileStorage) BatchUpdate(metrics []*metric.Metric) error {
	if err := s.memoryStorage.BatchUpdate(metrics); err != nil {
		return err
	}

	if s.interval != 0 {
		return nil
	}

	return s.sync()
}

func (s *FileStorage) Get(kind metric.Kind, name string) (*metric.Metric, error) {
	return s.memoryStorage.Get(kind, name)
}

func (s *FileStorage) List() ([]*metric.Metric, error) {
	return s.memoryStorage.List()
}

func (s *FileStorage) startSyncer() {
	go func() {
		sync := time.NewTicker(s.interval)
		defer sync.Stop()

		for {
			<-sync.C
			if err := s.sync(); err != nil {
				zlog.Logger.Errorf("sync metrics err=%w", err)
			}
		}
	}()
}

func (s *FileStorage) restore() error {
	data, err := os.ReadFile(s.filepath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}

		return fmt.Errorf("os readfile err=%w", err)
	}

	metrics := make([]*transport.Metric, 0)

	err = json.Unmarshal(data, &metrics)
	if err != nil {
		return fmt.Errorf("unmarshal err=%w", err)
	}

	gauges, counters := convertFromTransportMetrics(metrics)

	s.memoryStorage.CounterMetrics = counters
	s.memoryStorage.GaugeMetrics = gauges

	return nil
}

func (s *FileStorage) sync() error {
	data, err := s.serialize()
	if err != nil {
		return fmt.Errorf("serialize, err=%w", err)
	}

	err = os.WriteFile(s.filepath, data, 0644)
	if err != nil {
		return fmt.Errorf("write metrics to file=%s, err=%w", s.filepath, err)
	}

	return nil
}

func (s *FileStorage) serialize() ([]byte, error) {
	gauges := s.memoryStorage.GaugeMetrics.GetAll()
	counters := s.memoryStorage.CounterMetrics.GetAll()

	metrics := make([]*transport.Metric, 0, len(gauges)+len(counters))

	transportGauges, err := convertToTransportMetrics(gauges, metric.Gauge)
	if err != nil {
		return nil, err
	}

	transportCounters, err := convertToTransportMetrics(counters, metric.Counter)
	if err != nil {
		return nil, err
	}

	metrics = append(metrics, transportGauges...)
	metrics = append(metrics, transportCounters...)

	data, err := json.Marshal(metrics)
	if err != nil {
		return nil, fmt.Errorf("metrics marshal err=%w", err)
	}

	return data, nil
}

func (s *FileStorage) Stop() error {
	return nil
}

func convertToTransportMetrics[T int64 | float64](metrics map[string]T, kind metric.Kind) ([]*transport.Metric, error) {
	transportMetrics := make([]*transport.Metric, 0, len(metrics))

	for id, value := range metrics {
		m, err := transport.New(id, kind, value)
		if err != nil {
			return nil, fmt.Errorf("serializa id=%v kind=%v value=%v err=%w", id, kind, value, err)
		}

		transportMetrics = append(transportMetrics, m)
	}

	return transportMetrics, nil
}

func convertFromTransportMetrics(
	metrics []*transport.Metric,
) (
	*memorystorage.SyncStorage[float64],
	*memorystorage.SyncStorage[int64],
) {
	gauges := make(map[string]float64)
	counters := make(map[string]int64)

	for _, m := range metrics {
		switch m.Type {
		case metric.Gauge:
			gauges[m.ID] = *m.Value
		case metric.Counter:
			counters[m.ID] = *m.Delta
		default:
			zlog.Logger.Warnf("Unknown metric kind=%v", m.Type)
		}
	}

	syncGauges := memorystorage.NewSyncStorage[float64]()
	syncGauges.SetAll(gauges)

	syncCounters := memorystorage.NewSyncStorage[int64]()
	syncCounters.SetAll(counters)

	return syncGauges, syncCounters
}
