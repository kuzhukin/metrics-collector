package dbstorage

import (
	"github.com/kuzhukin/metrics-collector/internal/metric"
	"github.com/kuzhukin/metrics-collector/internal/server/storage"
)

const createGaugeMetricsTableQuery = `CREATE TABLE IF NOT EXISTS gauge_metrics (
	id text PRIMARY KEY,
	value double precision
);`

const createCounterMetricsTableQuery = `CREATE TABLE IF NOT EXISTS counter_metrics (
	id text PRIMARY KEY,
	value bigint
);`

func buildCreateMetricsTableQuery(kind metric.Kind) (string, error) {
	switch kind {
	case metric.Gauge:
		return createGaugeMetricsTableQuery, nil
	case metric.Counter:
		return createCounterMetricsTableQuery, nil
	default:
		return "", nil
	}
}

const updateGaugeMetricQuery = `INSERT INTO gauge_metrics (id, value) VALUES ($1, $2) ` +
	`ON CONFLICT (id) DO UPDATE SET value = excluded.value;`

const updateCounterMetricQuery = `INSERT INTO counter_metrics (id, value) VALUES ($1, $2) ` +
	`ON CONFLICT (id) DO UPDATE SET value = counter_metrics.value + excluded.value;`

func buildUpdateQuery(m *metric.Metric) (string, []interface{}, error) {
	switch m.Type {
	case metric.Gauge:
		return updateGaugeMetricQuery, []interface{}{m.ID, *m.Value}, nil
	case metric.Counter:
		return updateCounterMetricQuery, []interface{}{m.ID, *m.Delta}, nil
	default:
		return "", nil, storage.ErrUnknownKind
	}
}

func getdUpdateQueryByKind(k metric.Kind) (string, error) {
	switch k {
	case metric.Gauge:
		return updateGaugeMetricQuery, nil
	case metric.Counter:
		return updateCounterMetricQuery, nil
	default:
		return "", storage.ErrUnknownKind
	}
}

func prepareArgsForUpdate(m *metric.Metric) ([]interface{}, error) {
	switch m.Type {
	case metric.Gauge:
		return []interface{}{m.ID, *m.Value}, nil
	case metric.Counter:
		return []interface{}{m.ID, *m.Delta}, nil
	default:
		return nil, storage.ErrUnknownKind
	}
}

const getGaugeMetricQuery = `SELECT id, value FROM gauge_metrics WHERE id = $1;`
const getCounterMetricQuery = `SELECT id, value FROM counter_metrics WHERE id = $1;`

func buildGetQuery(id string, kind metric.Kind) (string, []interface{}, error) {
	switch kind {
	case metric.Gauge:
		return getGaugeMetricQuery, []interface{}{id}, nil
	case metric.Counter:
		return getCounterMetricQuery, []interface{}{id}, nil
	default:
		return "", nil, storage.ErrUnknownKind
	}
}

const getAllGaugeMetricsQuery = `SELECT * FROM gauge_metrics;`
const getAllCounterMetricsQuery = `SELECT * FROM counter_metrics;`

func buildGetAllQuery(kind metric.Kind) (string, error) {
	switch kind {
	case metric.Gauge:
		return getAllGaugeMetricsQuery, nil
	case metric.Counter:
		return getAllCounterMetricsQuery, nil
	default:
		return "", storage.ErrUnknownKind
	}
}
