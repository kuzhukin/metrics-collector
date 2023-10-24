package transport

import (
	"encoding/json"
	"errors"

	"github.com/kuzhukin/metrics-collector/internal/server/metric"
)

var ErrUnknownMetricType error = errors.New("unknown metric type")

type Metric struct {
	ID    string  `json:"id"`              // имя метрики
	Type  string  `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

func Serialize(id string, kind metric.Kind, value interface{}) ([]byte, error) {
	m := Metric{ID: id, Type: kind}

	switch kind {
	case metric.Counter:
		m.Delta = value.(int64)
	case metric.Gauge:
		m.Value = value.(float64)
	default:
		return nil, ErrUnknownMetricType
	}

	return json.Marshal(m)
}

func Desirialize(data []byte) (*Metric, error) {
	metric := &Metric{}
	err := json.Unmarshal(data, metric)

	return metric, err
}