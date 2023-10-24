package reporter

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/kuzhukin/metrics-collector/internal/log"
	"github.com/kuzhukin/metrics-collector/internal/server/metric"
	"github.com/kuzhukin/metrics-collector/internal/transport"
)

const updateEndpoint = "/update/"

//go:generate mockery --name=Reporter --filename=reporter.go --outpkg=mockreporter --output=mockreporter
type Reporter interface {
	Report(gaugeMetrics map[string]float64, counterMetrics map[string]int64)
}

type reporterImpl struct {
	updateURL string
}

func New(host string) Reporter {
	return &reporterImpl{
		updateURL: makeUpdateURL(host),
	}
}

func (r *reporterImpl) Report(gaugeMetrics map[string]float64, counterMetrics map[string]int64) {
	r.reportGauges(gaugeMetrics)
	r.reportCounters(counterMetrics)
}

func (r *reporterImpl) reportGauges(gaugeMetrics map[string]float64) {
	for name, value := range gaugeMetrics {
		if err := reportMetric(r.updateURL, name, metric.Gauge, value); err != nil {
			log.Logger.Warnf("report metric=%v, err=%s", name, err)
		}
	}
}

func (r *reporterImpl) reportCounters(counterMetrics map[string]int64) {
	for name, value := range counterMetrics {
		if err := reportMetric(r.updateURL, name, metric.Counter, value); err != nil {
			log.Logger.Warnf("report metric=%v, err=%s", name, err)
		}
	}
}

func reportMetric[T int64 | float64](URL string, id string, kind metric.Kind, value T) error {
	data, err := transport.Serialize(id, kind, value)
	if err != nil {
		return fmt.Errorf("metric serializa err=%s", err)
	}

	if err := doReport(URL, data); err != nil {
		return fmt.Errorf("do report, err=%w", err)
	}

	return nil
}

func doReport(URL string, data []byte) error {
	const maxTryingsNum = 5
	trying := 0

	var joinedError error

	for {
		if err := doOneReport(URL, data); err != nil {
			if trying < maxTryingsNum {
				trying++
				joinedError = errors.Join(joinedError, err)
				time.Sleep(time.Millisecond * 100)
				continue
			}

			return fmt.Errorf("http post URL=%v, err=%w", URL, joinedError)
		}

		return nil
	}
}

func doOneReport(URL string, data []byte) error {
	body := bytes.NewBuffer(data)
	resp, err := http.Post(URL, "application/json", body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("metrics report was failed with statusCode=%d", resp.StatusCode)
	}

	return nil
}

func makeUpdateURL(host string) string {
	return host + updateEndpoint
}
