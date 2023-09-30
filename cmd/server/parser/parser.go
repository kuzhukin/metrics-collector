package parser

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/kuzhukin/metrics-collector/internal/metric"
	"github.com/kuzhukin/metrics-collector/internal/shared"
)

const (
	kindIdx = iota
	nameIdx
	valueIdx
)

var ErrMetricNameIsNotFound error = errors.New("metric name isn't found")
var ErrBadMetricValue error = errors.New("bad metric value")
var ErrBadMetricKind error = errors.New("bad metric kind")

type RequestParser interface {
	Parse(request string) (*metric.Metric, error)
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

func (p *requestParserImpl) Parse(request string) (*metric.Metric, error) {
	metricRequest := cutUpdateEndpoint(request)

	params := strings.Split(metricRequest, "/")
	if len(params) != 3 || params[nameIdx] == "" {
		return nil, ErrMetricNameIsNotFound
	}

	v, err := p.parseValueByKind(params[kindIdx], params[valueIdx])
	if err != nil {
		return nil, err
	}

	return &metric.Metric{Kind: params[kindIdx], Name: params[nameIdx], Value: v}, nil
}

func cutUpdateEndpoint(urlPath string) string {
	metricPath, ok := strings.CutPrefix(urlPath, shared.UpdateEndpoint)
	if !ok {
		panic(fmt.Sprintf("bad metric format, prefix=%s, path=%s", shared.UpdateEndpoint, urlPath))
	}

	return metricPath
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
