package metric

type Value interface {
	Gauge() float64
	Counter() int64
}

type GaugeValue float64
type CounterValue float64

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
