package parser

import (
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/kuzhukin/metrics-collector/internal/metric"
	"github.com/kuzhukin/metrics-collector/internal/server/codec"
)

var _ RequestParser = &parserImpl{}
var _ BatchRequestParser = &parserImpl{}

type parserImpl struct {
}

func New() *parserImpl {
	return &parserImpl{}
}

func (p *parserImpl) Parse(r *http.Request) (*metric.Metric, error) {
	switch r.Header.Get("content-type") {
	case "application/json":
		return parseMetricByJSONBody(r)
	default:
		return parseMetricByURLParams(r)
	}
}

func (p *parserImpl) BatchParse(r *http.Request) ([]*metric.Metric, error) {
	switch r.Header.Get("content-type") {
	case "application/json":
		return parseBatchMetricByJSONBody(r)
	default:
		return nil, errors.New("incompatible content type for batch request")
	}
}

func parseMetricByURLParams(r *http.Request) (*metric.Metric, error) {
	kind := metric.Kind(chi.URLParam(r, "kind"))
	name := chi.URLParam(r, "name")
	valueStr := chi.URLParam(r, "value")

	if err := checkName(name); err != nil {
		return nil, err
	}

	if err := checkKind(kind); err != nil {
		return nil, err
	}

	m := &metric.Metric{Type: kind, ID: name}

	if valueStr != "" {
		delta, value, err := codec.Encode(kind, valueStr)
		if err != nil {
			return nil, err
		}

		m.Delta = delta
		m.Value = value
	}

	return m, nil
}

func parseMetricByJSONBody(r *http.Request) (*metric.Metric, error) {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading body from request, err=%w", err)
	}

	metricMsg, err := metric.Desirialize(data)
	if err != nil {
		return nil, fmt.Errorf("metric desirialization err=%w", err)
	}

	metric, err := parseMetricByJSONBodyImpl(metricMsg)
	if err != nil {
		return nil, fmt.Errorf("parse data=%v, err=%w", data, err)
	}

	return metric, nil
}

func parseBatchMetricByJSONBody(r *http.Request) ([]*metric.Metric, error) {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading body from request, err=%w", err)
	}

	batchMetric := metric.NewBatch()
	if err = batchMetric.Deserialize(data); err != nil {
		return nil, fmt.Errorf("batch metric desirialization err=%w", err)
	}

	metrics := make([]*metric.Metric, 0, batchMetric.Len())
	err = batchMetric.Foreach(func(nextMetric *metric.Metric) error {
		m, internalErr := parseMetricByJSONBodyImpl(nextMetric)
		if internalErr != nil {
			return fmt.Errorf("parse metric err=%w", err)
		}

		metrics = append(metrics, m)
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("parse batch data=%v, err=%w", data, err)
	}

	return metrics, nil
}

func parseMetricByJSONBodyImpl(metric *metric.Metric) (*metric.Metric, error) {
	if err := checkName(metric.ID); err != nil {
		return nil, fmt.Errorf("check name metric=%v, err=%w", metric, err)
	}

	if err := checkKind(metric.Type); err != nil {
		return nil, err
	}

	return metric, nil
}

func checkName(name string) error {
	if name == "" {
		return ErrMetricNameIsNotFound
	}

	return nil
}

func checkKind(kind metric.Kind) error {
	if kind != metric.Counter && kind != metric.Gauge {
		return ErrBadMetricKind
	}

	return nil
}
