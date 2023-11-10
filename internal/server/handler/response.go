package handler

import (
	"fmt"
	"net/http"

	"github.com/kuzhukin/metrics-collector/internal/metric"
	"github.com/kuzhukin/metrics-collector/internal/server/codec"
)

func response(w http.ResponseWriter, r *http.Request, metric *metric.Metric) error {
	switch r.Header.Get("Content-Type") {
	case "application/json":
		return responseJSON(w, r, metric)
	default:
		return responseTextPlain(w, r, metric)
	}
}

func responseJSON(w http.ResponseWriter, r *http.Request, m *metric.Metric) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	var data []byte
	var err error

	data, err = m.Serialize()
	if err != nil {
		return fmt.Errorf("serializa metric err=%w", err)
	}

	if _, err := w.Write(data); err != nil {
		return fmt.Errorf("write data err=%w", err)
	}

	return nil
}

func responseTextPlain(w http.ResponseWriter, r *http.Request, metric *metric.Metric) error {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)

	decodedValue := codec.DecodeValue(metric)
	if _, err := w.Write([]byte(decodedValue)); err != nil {
		return fmt.Errorf("write data, err=%w", err)
	}

	return nil
}
