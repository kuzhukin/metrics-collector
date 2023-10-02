package parser

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/kuzhukin/metrics-collector/internal/metric"
)

var _ RequestParser = &requestValueParserImpl{}

type requestValueParserImpl struct {
}

func NewValueRequestParser() RequestParser {
	return &requestValueParserImpl{}
}

func (p *requestValueParserImpl) Parse(r *http.Request) (*metric.Metric, error) {
	kind := chi.URLParam(r, "kind")
	name := chi.URLParam(r, "name")

	if name == "" {
		return nil, ErrMetricNameIsNotFound
	}

	return &metric.Metric{Kind: kind, Name: name}, nil
}
