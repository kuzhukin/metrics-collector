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

func ParseRequest(path string) (*metric.Metric, error) {
	metricPath := cutUpdateEndpoint(path)

	params := strings.Split(metricPath, "/")
	if len(params) != 3 || params[nameIdx] == "" {
		return nil, ErrMetricNameIsNotFound
	}

	return &metric.Metric{
		Kind:  params[kindIdx],
		Name:  params[nameIdx],
		Value: params[valueIdx],
	}, nil
}

func ParseGaugeValue(val string) (float64, error) {
	v, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return 0, errors.Join(ErrBadMetricValue, err)
	}

	return v, nil
}

func ParseCounterValue(val string) (int64, error) {
	v, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return 0, errors.Join(ErrBadMetricValue, err)
	}

	return v, nil
}

func cutUpdateEndpoint(urlPath string) string {
	metricPath, ok := strings.CutPrefix(urlPath, shared.UpdateEndpoint)
	if !ok {
		panic(fmt.Sprintf("bad metric format, prefix=%s, path=%s", shared.UpdateEndpoint, urlPath))
	}

	return metricPath
}
