package metric

import "errors"

type Value interface {
	Gauge() float64
	Counter() int64
}

type GaugeValue float64
type CounterValue int64

func (g GaugeValue) Gauge() float64 {
	return float64(g)
}

func (g GaugeValue) Counter() int64 {
	return int64(g)
}

func (g CounterValue) Gauge() float64 {
	return float64(g)
}

func (g CounterValue) Counter() int64 {
	return int64(g)
}

func NewValueByKind[T float64 | int64](kind Kind, value T) (Value, error) {
	switch kind {
	case Gauge:
		return GaugeValue(value), nil
	case Counter:
		return CounterValue(value), nil
	}

	return nil, errors.New("bad kind")
}
