package transport

import (
	"encoding/json"
	"errors"

	"github.com/kuzhukin/metrics-collector/internal/server/metric"
)

var ErrUnknownMetricType error = errors.New("unknown metric type")

type Metrics struct {
	ID    string   `json:"id"`              // имя метрики
	Type  string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

func Serialize(id string, kind metric.Kind, value interface{}) ([]byte, error) {
	m := Metrics{ID: id, Type: kind}

	switch kind {
	case metric.Counter:
		v := value.(int64)
		m.Delta = &v
	case metric.Gauge:
		v := value.(float64)
		m.Value = &v
	default:
		return nil, ErrUnknownMetricType
	}

	return json.Marshal(m)
}

func Desirialize(data []byte) (*Metrics, error) {
	metric := &Metrics{}
	err := json.Unmarshal(data, metric)

	return metric, err
}
