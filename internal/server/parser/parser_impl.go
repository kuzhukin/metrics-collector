package parser

import (
	"fmt"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/kuzhukin/metrics-collector/internal/server/codec"
	"github.com/kuzhukin/metrics-collector/internal/server/metric"
	"github.com/kuzhukin/metrics-collector/internal/transport"
)

var _ RequestParser = &updateRequestParserImpl{}

type updateRequestParserImpl struct {
}

func New() RequestParser {
	return &updateRequestParserImpl{}
}

func (p *updateRequestParserImpl) Parse(r *http.Request) (*metric.Metric, error) {
	switch r.Header.Get("content-type") {
	case "application/json":
		return parseMetricByJSONBody(r)
	default:
		return parseMetricByURLParams(r)
	}
}

func parseMetricByURLParams(r *http.Request) (*metric.Metric, error) {
	kind := chi.URLParam(r, "kind")
	name := chi.URLParam(r, "name")
	value := chi.URLParam(r, "value")

	if err := checkName(name); err != nil {
		return nil, err
	}

	if err := checkKind(kind); err != nil {
		return nil, err
	}

	var v metric.Value
	if value != "" {
		var err error
		v, err = codec.Encode(kind, value)
		if err != nil {
			return nil, err
		}
	}

	return &metric.Metric{Kind: kind, Name: name, Value: v}, nil
}

func parseMetricByJSONBody(r *http.Request) (*metric.Metric, error) {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading body from request, err=%w", err)
	}

	metricMsg, err := transport.Desirialize(data)
	if err != nil {
		return nil, fmt.Errorf("metric desirialization err=%w", err)
	}

	if err := checkName(metricMsg.ID); err != nil {
		return nil, fmt.Errorf("check name metric=%v, data=%v, err=%w", metricMsg, data, err)
	}

	if err := checkKind(metricMsg.Type); err != nil {
		return nil, err
	}

	var v metric.Value
	switch metricMsg.Type {
	case metric.Gauge:
		if metricMsg.Value != nil {
			v = metric.GaugeValue(*metricMsg.Value)
		}
	case metric.Counter:
		if metricMsg.Delta != nil {
			v = metric.CounterValue(*metricMsg.Delta)
		}
	}

	return &metric.Metric{Kind: metricMsg.Type, Name: metricMsg.ID, Value: v}, nil
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
