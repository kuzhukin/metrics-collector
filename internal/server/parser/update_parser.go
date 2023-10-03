package parser

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/kuzhukin/metrics-collector/internal/server/codec"
	"github.com/kuzhukin/metrics-collector/internal/server/metric"
)

var _ RequestParser = &requestUpdateParserImpl{}

type requestUpdateParserImpl struct {
}

func NewUpdateRequestParser() RequestParser {
	return &requestUpdateParserImpl{}
}

func (p *requestUpdateParserImpl) Parse(r *http.Request) (*metric.Metric, error) {
	kind := chi.URLParam(r, "kind")
	name := chi.URLParam(r, "name")
	value := chi.URLParam(r, "value")

	if name == "" {
		return nil, ErrMetricNameIsNotFound
	}

	if !isValidKind(kind) {
		return nil, ErrBadMetricKind
	}

	v, err := codec.Encode(kind, value)
	if err != nil {
		return nil, err
	}

	return &metric.Metric{Kind: kind, Name: name, Value: v}, nil
}

func isValidKind(kind metric.Kind) bool {
	return kind == metric.Counter || kind == metric.Gauge
}
