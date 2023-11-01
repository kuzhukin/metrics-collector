package transport

import (
	"encoding/json"
	"errors"

	"github.com/kuzhukin/metrics-collector/internal/server/metric"
)

var ErrUnknownMetricType error = errors.New("unknown metric type")

type Metric struct {
	ID    string   `json:"id"`              // имя метрики
	Type  string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

func New(id string, kind metric.Kind, value interface{}) (*Metric, error) {
	m := &Metric{ID: id, Type: kind}

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

	return m, nil
}

func (m *Metric) Serialize() ([]byte, error) {
	return json.Marshal(m)
}
func Serialize(id string, kind metric.Kind, value interface{}) ([]byte, error) {
	m, err := New(id, kind, value)
	if err != nil {
		return nil, err
	}

	return m.Serialize()
}

func Desirialize(data []byte) (*Metric, error) {
	metric := &Metric{}
	err := json.Unmarshal(data, metric)

	return metric, err
}
