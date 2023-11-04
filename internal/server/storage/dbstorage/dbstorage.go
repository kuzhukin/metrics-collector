package dbstorage

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/kuzhukin/metrics-collector/internal/server/metric"
	"github.com/kuzhukin/metrics-collector/internal/server/storage"
	"github.com/kuzhukin/metrics-collector/internal/zlog"
)

const pingTimeout = time.Second * 10

type DBStorage struct {
	db *sql.DB
}

var _ storage.Storage = &DBStorage{}

func StartNew(dataSourceName string) (*DBStorage, error) {
	db, err := sql.Open("pgx", dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("db conntection err=%w", err)
	}

	return &DBStorage{
		db: db,
	}, nil
}

func (s *DBStorage) Stop() error {
	s.db.Close()

	return nil
}

func (s *DBStorage) CheckConnection() bool {
	ctx, cancel := context.WithTimeout(context.Background(), pingTimeout)
	defer cancel()

	err := s.db.PingContext(ctx)

	if err != nil {
		zlog.Logger.Errorf("pint db err=%w", err)
	}

	return err == nil
}

func (s *DBStorage) Update(m *metric.Metric) error {
	return nil
}

func (s *DBStorage) Get(kind metric.Kind, name string) (*metric.Metric, error) {
	return nil, nil
}

func (s *DBStorage) List() []*metric.Metric {
	return nil
}
