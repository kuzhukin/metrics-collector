package reporter

import (
	"bytes"
	"compress/gzip"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
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
	tokenKey  []byte
}

func New(host string, tokenKey string) Reporter {
	var key []byte

	if tokenKey != "" {
		key = []byte(tokenKey)
	}

	return &reporterImpl{
		updateURL: makeUpdateURL(host),
		tokenKey:  key,
	}
}

func (r *reporterImpl) Report(gaugeMetrics map[string]float64, counterMetrics map[string]int64) {
	batch := prepareUpdate(gaugeMetrics, counterMetrics)
	if batch.Len() == 0 {
		return
	}

	if err := r.reportMetrics(batch); err != nil {
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

func (r *reporterImpl) reportMetrics(batch metric.MetricBatch) error {
	data, err := batch.Serialize()
	if err != nil {
		return fmt.Errorf("metric serializa err=%s", err)
	}

	if err := r.doReport(data); err != nil {
		return fmt.Errorf("do report, err=%w", err)
	}

	return nil
}

func (r *reporterImpl) doReport(data []byte) error {
	compressedData, err := compressData(data)
	if err != nil {
		return fmt.Errorf("compress data err=%w", err)
	}

	request, err := r.makeUpdateRequest(compressedData)
	if err != nil {
		return fmt.Errorf("make update request err=%w", err)
	}

	return doRequest(request)
}

func (r *reporterImpl) makeUpdateRequest(data []byte) (*http.Request, error) {
	req, err := http.NewRequest(http.MethodPost, r.updateURL, bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Encoding", "gzip")

	if r.tokenKey != nil {
		hash, err := r.hash(data)
		if err != nil {
			return nil, fmt.Errorf("hash data, err=%w", err)
		}

		req.Header.Set("HashSHA256", hex.EncodeToString(hash))
	}

	return req, nil
}

func (r *reporterImpl) hash(data []byte) ([]byte, error) {
	hasher := hmac.New(sha256.New, r.tokenKey)
	_, err := hasher.Write(data)
	if err != nil {
		return nil, fmt.Errorf("hasher write, err=%w", err)
	}

	return hasher.Sum(nil), nil
}

var tryingIntervals []time.Duration = []time.Duration{
	time.Millisecond * 1000,
	time.Millisecond * 3000,
	time.Millisecond * 5000,
}

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

	return fmt.Errorf("request trying limit exceeded, errs=%w", joinedError)
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
