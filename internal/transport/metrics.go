package transport

import (
	"encoding/json"
	"errors"
)

var ErrUnknownMetricType error = errors.New("unknown metric type")

type Metric struct {
	ID    string  `json:"id"`              // имя метрики
	Type  string  `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

func Serialize(id string, value interface{}) ([]byte, error) {
	metric := Metric{ID: id}

	switch v := value.(type) {
	case int64:
		metric.Type = "counter"
		metric.Delta = v
	case float64:
		metric.Type = "gauge"
		metric.Value = v
	default:
		return nil, ErrUnknownMetricType
	}

	return json.Marshal(metric)
}

func Desirialize(data []byte) (*Metric, error) {
	metric := &Metric{}
	err := json.Unmarshal(data, metric)

	return metric, err
}
