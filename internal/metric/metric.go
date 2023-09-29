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
	Value string
}

func (m *Metric) IsValid() bool {
	return m.Kind == Gauge || m.Kind == Counter
}
