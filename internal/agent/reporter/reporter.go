package reporter

import (
	"bytes"
	"compress/gzip"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/kuzhukin/metrics-collector/internal/server/metric"
	"github.com/kuzhukin/metrics-collector/internal/transport"
	"github.com/kuzhukin/metrics-collector/internal/zlog"
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
			zlog.Logger.Warnf("report metric=%v, err=%s", name, err)
		}
	}
}

func (r *reporterImpl) reportCounters(counterMetrics map[string]int64) {
	for name, value := range counterMetrics {
		if err := reportMetric(r.updateURL, name, metric.Counter, value); err != nil {
			zlog.Logger.Warnf("report metric=%v, err=%s", name, err)
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

	compressedData, err := compressData(data)
	if err != nil {
		return fmt.Errorf("compress data err=%w", err)
	}

	request, err := makeUpdateRequest(URL, compressedData)
	if err != nil {
		return fmt.Errorf("make update request err=%w", err)
	}

	return doRequest(request)
}

func makeUpdateRequest(URL string, data []byte) (*http.Request, error) {
	req, err := http.NewRequest(http.MethodPost, URL, bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Encoding", "gzip")

	return req, nil
}

func doRequest(req *http.Request) error {
	const maxTryingsNum = 5
	var joinedError error

	for trying := 0; trying < maxTryingsNum; trying++ {
		if resp, err := http.DefaultClient.Do(req); err != nil {
			if trying < maxTryingsNum {
				joinedError = errors.Join(joinedError, err)
				time.Sleep(time.Millisecond * 100)
			}
		} else {
			defer resp.Body.Close()
			if resp.StatusCode != http.StatusOK {
				return fmt.Errorf("metrics update request was failed with statusCode=%d", resp.StatusCode)
			}

			return nil
		}
	}

	return errors.New("request trying limit exceeded")
}

func makeUpdateURL(host string) string {
	return host + updateEndpoint
}

func compressData(data []byte) ([]byte, error) {
	b := &bytes.Buffer{}
	w := gzip.NewWriter(b)
	_, err := w.Write(data)
	if err != nil {
		return nil, err
	}

	err = w.Close()
	if err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}
