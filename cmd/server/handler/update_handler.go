package handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/kuzhukin/metrics-collector/cmd/server/metric"
	"github.com/kuzhukin/metrics-collector/cmd/server/parser"
	"github.com/kuzhukin/metrics-collector/cmd/server/shared"
	"github.com/kuzhukin/metrics-collector/cmd/server/storage"
)

var _ http.Handler = &UpdateHandler{}

type UpdateHandler struct {
	storage storage.Storage
}

func NewUpdateHandler(storage storage.Storage) *UpdateHandler {
	return &UpdateHandler{
		storage: storage,
	}
}

func (u *UpdateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		fmt.Printf("Endpoint %s supports only POST method\n", shared.UpdateEndpoint)
		w.WriteHeader(http.StatusMethodNotAllowed)

		return
	}

	metric, err := parser.ParseRequest(r.URL.Path)
	if err != nil {
		if errors.Is(err, parser.ErrMetricNameIsNotFound) {
			fmt.Printf("Metric name isn't found path=%s, err=%s\n", r.URL.Path, err)
			w.WriteHeader(http.StatusNotFound)

			return
		} else {
			fmt.Printf("Parse request error=%s\n", err)
			w.WriteHeader(http.StatusBadRequest)

			return
		}
	}

	if !metric.IsValid() {
		fmt.Printf("Bad metric=%v\n", metric)
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	if err := u.updateMetric(metric); err != nil {
		if errors.Is(err, parser.ErrBadMetricValue) {
			fmt.Printf("Bad value of metric=%v\n", metric)
			w.WriteHeader(http.StatusBadRequest)

			return
		} else {
			fmt.Printf("Metrics=%v updating err=%s\n", metric, err)
			w.WriteHeader(http.StatusInternalServerError)

			return
		}
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
}

func (u *UpdateHandler) updateMetric(m *metric.Metric) error {
	switch m.Kind {
	case metric.Gauge:
		return u.updateGauge(m.Name, m.Value)
	case metric.Counter:
		return u.updateCounter(m.Name, m.Value)
	default:
		return fmt.Errorf("unknown metric kind=%s", m.Kind)
	}
}

func (u *UpdateHandler) updateGauge(name string, valueStr string) error {
	value, err := parser.ParseGaugeValue(valueStr)
	if err != nil {
		return fmt.Errorf("parse gauge value=%s, err=%w", valueStr, err)
	}

	if err := u.storage.UpdateGauge(name, value); err != nil {
		return fmt.Errorf("update gauge name=%s value=%v, err=%w", name, value, err)
	}

	return nil
}

func (u *UpdateHandler) updateCounter(name string, valueStr string) error {
	value, err := parser.ParseCounterValue(valueStr)
	if err != nil {
		return fmt.Errorf("parse counter value=%s, err=%w", valueStr, err)
	}

	if err := u.storage.UpdateCounter(name, value); err != nil {
		return fmt.Errorf("update counter name=%s value=%v, err=%w", name, value, err)
	}

	return nil
}
