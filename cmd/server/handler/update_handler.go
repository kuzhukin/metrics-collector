package handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/kuzhukin/metrics-collector/cmd/server/parser"
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

	metric, err := u.parser.Parse(r)
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

	if err := u.storage.Update(metric); err != nil {
		fmt.Printf("Metrics=%v updating err=%s\n", metric, err)
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
}
