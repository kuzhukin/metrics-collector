package metric

import "encoding/json"

type MetricBatch struct {
	metrics []*Metric
}

func NewBatch() MetricBatch {
	return MetricBatch{
		metrics: make([]*Metric, 0),
	}
}

func (b *MetricBatch) Add(m *Metric) {
	b.metrics = append(b.metrics, m)
}

func (b *MetricBatch) Len() int {
	return len(b.metrics)
}

func (b *MetricBatch) Serialize() ([]byte, error) {
	return json.Marshal(b.metrics)
}

func (b *MetricBatch) Deserialize(data []byte) error {
	return json.Unmarshal(data, &b.metrics)
}

func (b *MetricBatch) Foreach(fn func(m *Metric) error) error {
	for _, m := range b.metrics {
		if err := fn(m); err != nil {
			return err
		}
	}

	return nil
}
