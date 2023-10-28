package filestorage

import (
	"time"

	"github.com/kuzhukin/metrics-collector/internal/server/config"
	"github.com/kuzhukin/metrics-collector/internal/server/metric"
	"github.com/kuzhukin/metrics-collector/internal/server/storage"
	"github.com/kuzhukin/metrics-collector/internal/server/storage/memorystorage"
	"github.com/kuzhukin/metrics-collector/internal/transport"
)

var _ storage.Storage = &FileStorage{}

type FileStorage struct {
	memorystorage.MemoryStorage

	filepath string
	interval time.Duration
	done     chan struct{}
}

type StorageEncoder interface {
	Encode() map[string]transport.Metric
}

func New(config config.StorageConfig) (*FileStorage, error) {
	storage := &FileStorage{
		MemoryStorage: memorystorage.MemoryStorage{
			GaugeMetrics:   memorystorage.NewSyncStorage[float64](),
			CounterMetrics: memorystorage.NewSyncStorage[int64](),
		},

		filepath: config.FilePath,
		interval: time.Second * time.Duration(config.Interval),
	}

	if config.Restore {
		storage.restore()
	}

	if storage.interval > 0 {
		storage.startSyncer()
	}

	return storage, nil
}

func (s *FileStorage) restore() error {
	// TODO:
	return nil
}

func (s *FileStorage) startSyncer() {
	s.done = make(chan struct{})

	go func() {
		sync := time.NewTicker(s.interval)
		defer sync.Stop()

		select {
		case <-sync.C:
			s.sync()
		case <-s.done:
			return
		}
	}()
}

func (s *FileStorage) sync() error {
	// TODO:
	return nil
}

func (s *FileStorage) Update(m *metric.Metric) error {
	return nil
}

func (s *FileStorage) Get(kind metric.Kind, name string) (*metric.Metric, error) {
	return nil, nil
}

func (s *FileStorage) List() []*metric.Metric {
	return nil
}
