package reporter

import (
	"bytes"
	"fmt"
	"net/http"
	"strconv"

	"github.com/kuzhukin/metrics-collector/internal/server/metric"
)

const updateEndpoint = "/update/"

//go:generate mockgen -source=reporter.go -destination=mockreporter/mock.go -package=mockreporter
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
		encodedValue := strconv.FormatFloat(value, 'G', -1, 64)
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
	resp, err := http.Post(request, "text/plain", body)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("request=%s was failed with statusCode=%d\n", request, resp.StatusCode)
	}

	return nil
}

func makeUpdateRequest(hostport string, kind metric.Kind, name string, value string) string {
	return hostport + updateEndpoint + kind + "/" + name + "/" + value
}
