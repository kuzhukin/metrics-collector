package reporter

import (
	"bytes"
	"net/http"
	"strconv"

	"github.com/kuzhukin/metrics-collector/internal/log"
	"github.com/kuzhukin/metrics-collector/internal/server/metric"
)

const updateEndpoint = "/update/"

//go:generate mockery --name=Reporter --filename=reporter.go --outpkg=mockreporter --output=mockreporter
type Reporter interface {
	Report(gaugeMetrics map[string]float64, counterMetrics map[string]int64)
}

type reporterImpl struct {
	URL string
}

func New(URL string) Reporter {
	return &reporterImpl{
		URL: URL,
	}
}

func (r *reporterImpl) Report(gaugeMetrics map[string]float64, counterMetrics map[string]int64) {
	r.reportGauge(gaugeMetrics)
	r.reportCounter(counterMetrics)
}

func (r *reporterImpl) reportGauge(gaugeMetrics map[string]float64) {
	for name, value := range gaugeMetrics {
		encodedValue := strconv.FormatFloat(value, 'G', -1, 64)
		request := makeUpdateRequest(r.URL, metric.Gauge, name, encodedValue)

		if err := doReport(request); err != nil {
			log.Logger.Warnf("Do gauge metric report=%s, err=%s\n", request, err)
		}
	}
}

func (r *reporterImpl) reportCounter(counterMetrics map[string]int64) {
	for name, value := range counterMetrics {
		encodedValue := strconv.FormatInt(value, 10)
		request := makeUpdateRequest(r.URL, metric.Counter, name, encodedValue)

		if err := doReport(request); err != nil {
			log.Logger.Warnf("Do counters metric report=%s, err=%s\n", request, err)
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
		log.Logger.Warnf("request=%s was failed with statusCode=%d\n", request, resp.StatusCode)
	}

	return nil
}

func makeUpdateRequest(URL string, kind metric.Kind, name string, value string) string {
	return URL + updateEndpoint + kind + "/" + name + "/" + value
}
