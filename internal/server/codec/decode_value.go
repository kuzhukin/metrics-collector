package codec

import (
	"strconv"

	"github.com/kuzhukin/metrics-collector/internal/metric"
)

type dencodeFunc = func(m *metric.Metric) string

var valueDecoders = map[metric.Kind]dencodeFunc{
	metric.Gauge:   dencodeGauge,
	metric.Counter: dencodeCounter,
}

func DecodeValue(m *metric.Metric) string {
	return valueDecoders[m.Type](m)
}

func dencodeGauge(m *metric.Metric) string {
	return strconv.FormatFloat(*m.Value, 'G', -1, 64)
}

func dencodeCounter(m *metric.Metric) string {
	return strconv.FormatInt(*m.Delta, 10)
}
