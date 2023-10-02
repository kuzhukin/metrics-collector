package codec

import (
	"errors"
	"strconv"

	"github.com/kuzhukin/metrics-collector/internal/metric"
)

var ErrBadMetricValue error = errors.New("bad metric value")

type encodeFunc = func(val string) (metric.Value, error)

var valueEncoders = map[metric.Kind]encodeFunc{
	metric.Gauge:   encodeGauge,
	metric.Counter: encodeCounter,
}

func Encode(kind metric.Kind, value string) (metric.Value, error) {
	return valueEncoders[kind](value)
}

func encodeGauge(val string) (metric.Value, error) {
	v, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return nil, errors.Join(ErrBadMetricValue, err)
	}

	return metric.GaugeValue(v), nil
}

func encodeCounter(val string) (metric.Value, error) {
	v, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return nil, errors.Join(ErrBadMetricValue, err)
	}

	return metric.CounterValue(v), nil
}
