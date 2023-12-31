package dbstorage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/kuzhukin/metrics-collector/internal/metric"
	"github.com/kuzhukin/metrics-collector/internal/server/storage"
	"github.com/kuzhukin/metrics-collector/internal/zlog"
)

const (
	pingTimeout            = time.Second * 10
	createTablesTimeout    = time.Second * 10
	updateMetricTimeout    = time.Second * 10
	updateAllMetricTimeout = time.Second * 60
	getMetricTimeout       = time.Second * 10
	getAllMetricsTimeout   = time.Second * 60
)

var tryingIntervals = []time.Duration{time.Second * 1, time.Second * 3, time.Second * 5}

var (
	compatibleMetricKinds = []metric.Kind{
		metric.Counter,
		metric.Gauge,
	}
)

type DBStorage struct {
	db *sql.DB
}

var _ storage.Storage = &DBStorage{}

func StartNew(dataSourceName string) (*DBStorage, error) {
	db, err := sql.Open("pgx", dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("db conntection err=%w", err)
	}

	storage := &DBStorage{db: db}

	if err := storage.createTables(); err != nil {
		return nil, err
	}

	return storage, nil
}

func (s *DBStorage) createTables() error {
	var err error

	for _, kind := range compatibleMetricKinds {
		err = errors.Join(err, s.createTableForKind(kind))
	}

	return err
}

func (s *DBStorage) createTableForKind(kind metric.Kind) error {
	query, err := buildCreateMetricsTableQuery(kind)
	if err != nil {
		return fmt.Errorf("build create table query for kind=%v, err=%w", kind, err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), createTablesTimeout)
	defer cancel()

	_, err = s.db.ExecContext(ctx, query)
	if err != nil {
		return fmt.Errorf("exec create table query, err=%w", err)
	}

	return err
}

func (s *DBStorage) Stop() error {
	s.db.Close()

	return nil
}

func (s *DBStorage) CheckConnection(ctx context.Context) bool {
	ctx, cancel := context.WithTimeout(ctx, pingTimeout)
	defer cancel()

	err := s.db.PingContext(ctx)

	if err != nil {
		zlog.Logger.Errorf("pint db err=%w", err)
	}

	return err == nil
}

func (s *DBStorage) Update(ctx context.Context, m *metric.Metric) error {
	query, args, err := buildUpdateQuery(m)
	if err != nil {
		return fmt.Errorf("build update query for metric=%v, err=%w", m, err)
	}

	execFunc := func() (*sql.Result, error) {
		ctx, cancel := context.WithTimeout(ctx, updateMetricTimeout)
		defer cancel()

		res, err := s.db.ExecContext(ctx, query, args...)
		if err != nil {
			return nil, fmt.Errorf("exec metric=%v update err=%w", m, err)
		}

		return &res, nil
	}

	_, err = doQuery(execFunc)
	if err != nil {
		return fmt.Errorf("do query, err=%w", err)
	}

	return nil
}

func (s *DBStorage) BatchUpdate(ctx context.Context, metrics []*metric.Metric) error {
	groupedMetrics := groupMetricsByKind(metrics)

	query := func() (*struct{}, error) {
		if err := s.updateMetrics(ctx, groupedMetrics); err != nil {
			return nil, fmt.Errorf("update metrics, err=%w", err)
		}

		return nil, nil
	}

	if _, err := doQuery(query); err != nil {
		return fmt.Errorf("do query, error")
	}

	return nil
}

func (s *DBStorage) updateMetrics(ctx context.Context, metricsByKind map[metric.Kind][]*metric.Metric) error {
	ctx, cancel := context.WithTimeout(ctx, updateAllMetricTimeout)
	defer cancel()

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx, err=%w", err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	statements := make([]*sql.Stmt, 0)
	defer func() {
		for _, st := range statements {
			st.Close()
		}
	}()

	for kind, metrics := range metricsByKind {
		query, err := getdUpdateQueryByKind(kind)
		if err != nil {
			return fmt.Errorf("get update query by kind")
		}

		stmt, err := tx.PrepareContext(ctx, query)
		if err != nil {
			return err
		}
		statements = append(statements, stmt)

		for _, m := range metrics {
			args, err := prepareArgsForUpdate(m)
			if err != nil {
				return fmt.Errorf("prepare args for metric name=%s, kind=%s, err=%w", m.ID, m.Type, err)
			}

			_, err = stmt.ExecContext(ctx, args...)
			if err != nil {
				return fmt.Errorf("stmt exec, err=%w", err)
			}

		}
	}

	return tx.Commit()
}

func groupMetricsByKind(metrics []*metric.Metric) map[metric.Kind][]*metric.Metric {
	grouped := make(map[metric.Kind][]*metric.Metric)

	for _, m := range metrics {
		kindMetrics, ok := grouped[m.Type]
		if !ok {
			kindMetrics = []*metric.Metric{m}
		} else {
			kindMetrics = append(kindMetrics, m)
		}

		grouped[m.Type] = kindMetrics
	}

	return grouped
}

func (s *DBStorage) Get(ctx context.Context, kind metric.Kind, name string) (*metric.Metric, error) {
	query, args, err := buildGetQuery(name, kind)
	if err != nil {
		return nil, fmt.Errorf("build get query, err=%w", err)
	}

	queryFunc := func() (*sql.Rows, error) {
		ctx, cancel := context.WithTimeout(ctx, getMetricTimeout)
		defer cancel()

		rows, err := s.db.QueryContext(ctx, query, args...)
		if err != nil {
			return nil, fmt.Errorf("query metric, err=%w", err)
		}

		return rows, nil
	}

	rows, err := doQuery(queryFunc)
	if err != nil {
		return nil, fmt.Errorf("do query, err=%w", err)
	}
	defer rows.Close()

	if rows.Next() {
		parser, err := makeParserForKind(kind)
		if err != nil {
			return nil, err
		}

		return parser(rows)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return nil, storage.ErrUnknownMetric
}

func (s *DBStorage) List(ctx context.Context) ([]*metric.Metric, error) {
	acc := make([]*metric.Metric, 0)

	for _, kind := range compatibleMetricKinds {
		metrics, err := s.getAll(ctx, kind)
		if err != nil {
			return nil, err
		}

		acc = append(acc, metrics...)
	}

	return acc, nil
}

func (s *DBStorage) getAll(ctx context.Context, kind metric.Kind) ([]*metric.Metric, error) {
	query, err := buildGetAllQuery(kind)
	if err != nil {
		return nil, fmt.Errorf("build get all query for kind=%v, err=%w", kind, err)
	}

	queryFunc := func() (*sql.Rows, error) {
		ctx, cancel := context.WithTimeout(ctx, getAllMetricsTimeout)
		defer cancel()

		rows, err := s.db.QueryContext(ctx, query)
		if err != nil {
			return nil, fmt.Errorf("query all metrics with kind=%s, err=%w", kind, err)
		}

		return rows, nil
	}

	rows, err := doQuery(queryFunc)
	if err != nil {
		return nil, fmt.Errorf("do query, err=%w", err)
	}
	defer rows.Close()

	parser, err := makeParserForKind(kind)
	if err != nil {
		return nil, err
	}

	metrics := make([]*metric.Metric, 0)

	for rows.Next() {
		m, err := parser(rows)
		if err != nil {
			return nil, fmt.Errorf("parse rows, err=%w", err)
		}

		metrics = append(metrics, m)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return metrics, nil
}

func makeParserForKind(kind metric.Kind) (func(rows *sql.Rows) (*metric.Metric, error), error) {
	switch kind {
	case metric.Gauge:
		return func(innerRows *sql.Rows) (*metric.Metric, error) {
			name, value, err := parse[float64](innerRows)
			if err != nil {
				return nil, err
			}

			return &metric.Metric{ID: name, Type: metric.Gauge, Value: &value}, nil
		}, nil
	case metric.Counter:
		return func(innerRows *sql.Rows) (*metric.Metric, error) {
			name, value, err := parse[int64](innerRows)
			if err != nil {
				return nil, err
			}

			return &metric.Metric{ID: name, Type: metric.Counter, Delta: &value}, nil
		}, nil
	default:
		return nil, storage.ErrUnknownKind
	}
}

func parse[T int64 | float64](rows *sql.Rows) (string, T, error) {
	name := ""
	value := T(0)

	if err := rows.Scan(&name, &value); err != nil {
		return "", value, fmt.Errorf("scan metric row, err=%w", err)
	}

	return name, value, nil
}

func doQuery[T any](queryFunc func() (*T, error)) (*T, error) {
	var commonErr error
	max := len(tryingIntervals)

	for trying := 0; trying <= max; trying++ {
		rows, err := queryFunc()
		if err != nil {
			commonErr = errors.Join(commonErr, err)

			if trying < max && isRetriableError(err) {
				time.Sleep(tryingIntervals[trying])
				continue
			}

			return nil, commonErr
		}

		return rows, nil
	}

	return nil, commonErr
}

func isRetriableError(err error) bool {
	var pgErr *pgconn.PgError

	return errors.As(err, &pgErr) && (pgerrcode.IsIntegrityConstraintViolation(pgErr.Code) || pgerrcode.IsConnectionException(pgErr.Code))
}
