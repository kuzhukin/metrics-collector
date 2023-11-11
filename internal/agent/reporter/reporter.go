package reporter

import (
	"bytes"
	"compress/gzip"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/kuzhukin/metrics-collector/internal/metric"
	"github.com/kuzhukin/metrics-collector/internal/zlog"
)

const batchUpdateEndpoint = "/updates/"

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
	batch := prepareUpdate(gaugeMetrics, counterMetrics)
	if batch.Len() == 0 {
		return
	}

	if err := reportMetrics(r.updateURL, batch); err != nil {
		zlog.Logger.Errorf("report metrics err=%s", err)
	}
}

func prepareUpdate(gaugeMetrics map[string]float64, counterMetrics map[string]int64) metric.MetricBatch {
	batch := metric.NewBatch()

	batch = prepare(gaugeMetrics, metric.Gauge, batch)
	batch = prepare(counterMetrics, metric.Counter, batch)

	return batch
}

func prepare[T int64 | float64](metrics map[string]T, kind metric.Kind, acc metric.MetricBatch) metric.MetricBatch {
	for name, value := range metrics {
		m, err := metric.New(name, kind, value)
		if err != nil {
			zlog.Logger.Warnf("report metric=%v, err=%s", name, err)
			continue
		}

		acc.Add(m)
	}

	return acc
}

func reportMetrics(URL string, batch metric.MetricBatch) error {
	data, err := batch.Serialize()
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

var tryingIntervals []time.Duration = []time.Duration{time.Second * 1, time.Second * 3, time.Second * 5}

func doRequest(req *http.Request) error {

	var joinedError error
	maxTryingsNum := len(tryingIntervals)

	for trying := 0; trying <= maxTryingsNum; trying++ {
		if resp, err := http.DefaultClient.Do(req); err != nil {
			if trying < maxTryingsNum {
				joinedError = errors.Join(joinedError, err)
				time.Sleep(tryingIntervals[trying])
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
	return host + batchUpdateEndpoint
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
