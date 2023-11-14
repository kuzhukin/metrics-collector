package codec

import (
	"errors"
	"strconv"

	"github.com/kuzhukin/metrics-collector/internal/metric"
)

var ErrBadMetricValue error = errors.New("bad metric value")

type encodeFunc = func(val string) (*int64, *float64, error)

var valueEncoders = map[metric.Kind]encodeFunc{
	metric.Gauge:   encodeGauge,
	metric.Counter: encodeCounter,
}

func Encode(kind metric.Kind, value string) (*int64, *float64, error) {
	return valueEncoders[kind](value)
}

func encodeGauge(val string) (*int64, *float64, error) {
	v, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return nil, nil, errors.Join(ErrBadMetricValue, err)
	}

	return nil, &v, nil
}

func encodeCounter(val string) (*int64, *float64, error) {
	v, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return nil, nil, errors.Join(ErrBadMetricValue, err)
	}

	return &v, nil, nil
}
