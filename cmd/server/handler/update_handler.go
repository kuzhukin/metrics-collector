package handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/kuzhukin/metrics-collector/cmd/server/parser"
	"github.com/kuzhukin/metrics-collector/internal/metric"
	"github.com/kuzhukin/metrics-collector/internal/shared"
	"github.com/kuzhukin/metrics-collector/internal/storage"
)

var _ http.Handler = &UpdateHandler{}

type UpdateHandler struct {
	storage storage.Storage
	parser  parser.RequestParser
}

func NewUpdateHandler(storage storage.Storage, parser parser.RequestParser) *UpdateHandler {
	return &UpdateHandler{
		storage: storage,
		parser:  parser,
	}
}

func (u *UpdateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		fmt.Printf("Endpoint %s supports only POST method\n", shared.UpdateEndpoint)
		w.WriteHeader(http.StatusMethodNotAllowed)

		return
	}

	metric, err := u.parser.Parse(r.URL.Path)
	if err != nil {
		if errors.Is(err, parser.ErrMetricNameIsNotFound) {
			fmt.Printf("Metric name isn't found path=%s, err=%s\n", r.URL.Path, err)
			w.WriteHeader(http.StatusNotFound)
		} else if errors.Is(err, parser.ErrBadMetricValue) {
			fmt.Printf("Bad metric value path=%s, err=%s\n", r.URL.Path, err)
			w.WriteHeader(http.StatusBadRequest)
		} else if errors.Is(err, parser.ErrBadMetricKind) {
			fmt.Printf("Bad metric kind path=%s, err=%s\n", r.URL.Path, err)
			w.WriteHeader(http.StatusBadRequest)
		} else {
			fmt.Printf("Parse request path=%s, err=%s\n", r.URL.Path, err)
			w.WriteHeader(http.StatusBadRequest)
		}

		return
	}

	if err := u.updateMetric(metric); err != nil {
		fmt.Printf("Metrics=%v updating err=%s\n", metric, err)
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
}

func (u *UpdateHandler) updateMetric(m *metric.Metric) error {
	switch m.Kind {
	case metric.Gauge:
		return u.updateGauge(m.Name, m.Value.Gauge())
	case metric.Counter:
		return u.updateCounter(m.Name, m.Value.Counter())
	default:
		return fmt.Errorf("doesn't have update handle func for kind=%s", m.Kind)
	}
}

func (u *UpdateHandler) updateGauge(name string, value float64) error {
	if err := u.storage.UpdateGauge(name, value); err != nil {
		return fmt.Errorf("update gauge name=%s value=%v, err=%w", name, value, err)
	}

	return nil
}

func (u *UpdateHandler) updateCounter(name string, value int64) error {
	if err := u.storage.UpdateCounter(name, value); err != nil {
		return fmt.Errorf("update counter name=%s value=%v, err=%w", name, value, err)
	}

	return nil
}
