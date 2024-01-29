package metric

import (
	"encoding/json"
	"errors"
)

type Kind string

const (
	Gauge   Kind = "gauge"
	Counter Kind = "counter"
)

var ErrUnknownMetricType error = errors.New("unknown metric type")

// a model of single metric
type Metric struct {
	ID    string   `json:"id"`
	Type  Kind     `json:"type"`
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
}

func New(id string, kind Kind, value interface{}) (*Metric, error) {
	m := &Metric{ID: id, Type: kind}

	switch kind {
	case Counter:
		v := value.(int64)
		m.Delta = &v
	case Gauge:
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

func Serialize(id string, kind Kind, value interface{}) ([]byte, error) {
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
