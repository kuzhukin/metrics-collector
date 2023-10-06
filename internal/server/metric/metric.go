package metric

type Kind = string

const (
	// is float metric
	Gauge Kind = "gauge"
	// is integer metric
	Counter Kind = "counter"
)

type Metric struct {
	Kind  Kind
	Name  string
	Value Value
}

func NewMetric(kind Kind, name string, value Value) *Metric {
	return &Metric{
		Kind:  kind,
		Name:  name,
		Value: value,
	}
}
