package parser

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/kuzhukin/metrics-collector/internal/metric"
)

var ErrMetricNameIsNotFound error = errors.New("metric name isn't found")
var ErrBadMetricValue error = errors.New("bad metric value")
var ErrBadMetricKind error = errors.New("bad metric kind")

//go:generate mockgen -source=parser.go -destination=mockparser/mock.go -package=mockparser
type RequestParser interface {
	Parse(r *http.Request) (*metric.Metric, error)
}

type requestParserImpl struct {
	valueParsers map[metric.Kind]valueParser
}

func NewRequestParser() RequestParser {
	return &requestParserImpl{
		valueParsers: map[metric.Kind]valueParser{
			metric.Gauge:   parseGaugeValue,
			metric.Counter: parseCounterValue,
		},
	}
}

func (p *requestParserImpl) Parse(r *http.Request) (*metric.Metric, error) {
	kind := chi.URLParam(r, "kind")
	name := chi.URLParam(r, "name")
	value := chi.URLParam(r, "value")

	if name == "" {
		return nil, ErrMetricNameIsNotFound
	}

	v, err := p.parseValueByKind(kind, value)
	if err != nil {
		return nil, err
	}

	return &metric.Metric{Kind: kind, Name: name, Value: v}, nil
}

type valueParser = func(val string) (metric.Value, error)

func (p *requestParserImpl) parseValueByKind(kind metric.Kind, value string) (metric.Value, error) {
	parser, ok := p.valueParsers[kind]
	if !ok {
		return nil, ErrBadMetricKind
	}

	return parser(value)
}

func parseGaugeValue(val string) (metric.Value, error) {
	v, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return nil, errors.Join(ErrBadMetricValue, err)
	}

	return metric.GaugeValue(v), nil
}

func parseCounterValue(val string) (metric.Value, error) {
	v, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return nil, errors.Join(ErrBadMetricValue, err)
	}

	return metric.CounterValue(v), nil
}
