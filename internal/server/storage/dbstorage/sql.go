package dbstorage

import (
	"github.com/kuzhukin/metrics-collector/internal/server/metric"
	"github.com/kuzhukin/metrics-collector/internal/server/storage"
)

const createGaugeMetricsTableQuery = `CREATE TABLE iF NOT EXISTS gauge_metrics (
	"id" text,
	"value" double precision,
	PRIMARY KEY "id"
);`

const createCounterMetricsTableQuery = `CREATE TABLE iF NOT EXISTS counter_metrics (
	"id" text,
	"value" bigint,
	PRIMARY KEY "id"
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

const updateGaugeMetricQuery = `INSERT INTO gauge_metrics ("id", "value") VALUES ($1, $2) ` +
	`ON CONFLICT ("id", "value") metrics.value = excluded.value;`

const updateCounterMetricQuery = `INSERT INTO counter_metrics ("id", "value") VALUES ($1, $2) ` +
	`ON CONFLICT ("id", "value") metrics.value = metrics.value + excluded.value;`

func buildUpdateQuery(m *metric.Metric) (string, []interface{}, error) {
	switch m.Kind {
	case metric.Gauge:
		return updateGaugeMetricQuery, []interface{}{m.Name, m.Value.Gauge()}, nil
	case metric.Counter:
		return updateCounterMetricQuery, []interface{}{m.Name, m.Value.Counter()}, nil
	default:
		return "", nil, storage.ErrUnknownKind
	}
}

const getGaugeMetricQuery = `SELECT name, value FROM gauge_metrics WHERE name = $1;`
const getCounterMetricQuery = `SELECT name, value FROM counter_metrics WHERE name = $1;`

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
