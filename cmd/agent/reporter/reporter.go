package reporter

import (
	"bytes"
	"fmt"
	"net/http"
	"strconv"

	"github.com/kuzhukin/metrics-collector/internal/metric"
	"github.com/kuzhukin/metrics-collector/internal/shared"
)

type Reporter interface {
	Report(gaugeMetrics map[string]float64, counterMetrics map[string]int64)
}

type reporterImpl struct {
	hostport string
}

func New(hostport string) Reporter {
	return &reporterImpl{
		hostport: hostport,
	}
}

func (r *reporterImpl) Report(gaugeMetrics map[string]float64, counterMetrics map[string]int64) {
	r.reportGauge(gaugeMetrics)
	r.reportCounter(counterMetrics)
}

func (r *reporterImpl) reportGauge(gaugeMetrics map[string]float64) {
	for name, value := range gaugeMetrics {
		encodedValue := strconv.FormatFloat(value, 'f', 4, 64)
		request := makeUpdateRequest(r.hostport, metric.Gauge, name, encodedValue)

		if err := doReport(request); err != nil {
			fmt.Printf("Do gauge metric report=%s, err=%s\n", request, err)
		}
	}
}

func (r *reporterImpl) reportCounter(counterMetrics map[string]int64) {
	for name, value := range counterMetrics {
		encodedValue := strconv.FormatInt(value, 10)
		request := makeUpdateRequest(r.hostport, metric.Counter, name, encodedValue)

		if err := doReport(request); err != nil {
			fmt.Printf("Do counters metric report=%s, err=%s\n", request, err)
		}
	}
}

func doReport(request string) error {
	body := bytes.NewBufferString("")
	_, err := http.Post(request, "text/plain", body)

	return err
}

func makeUpdateRequest(hostport string, kind metric.Kind, name string, value string) string {
	return hostport + shared.UpdateEndpoint + kind + "/" + name + "/" + value
}
