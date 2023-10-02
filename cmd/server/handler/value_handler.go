package handler

import (
	"fmt"
	"io"
	"net/http"

	"github.com/kuzhukin/metrics-collector/cmd/server/codec"
	"github.com/kuzhukin/metrics-collector/cmd/server/parser"
	"github.com/kuzhukin/metrics-collector/internal/shared"
	"github.com/kuzhukin/metrics-collector/internal/storage"
)

var _ http.Handler = &ValueHandler{}

type ValueHandler struct {
	storage storage.Storage
	parser  parser.RequestParser
}

func NewValueHandler(storage storage.Storage, parser parser.RequestParser) *ValueHandler {
	return &ValueHandler{
		storage: storage,
		parser:  parser,
	}
}

func (u *ValueHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		fmt.Printf("Endpoint %s supports only GET method\n", shared.ValueEndpoint)
		w.WriteHeader(http.StatusMethodNotAllowed)

		return
	}

	metric, err := u.parser.Parse(r)
	if err != nil {
		fmt.Printf("Parse request path=%s, err=%s\n", r.URL.Path, err)
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	metric, err = u.storage.Get(metric.Kind, metric.Name)
	if err != nil {
		fmt.Printf("Metrics=%v updating err=%s\n", metric, err)
		w.WriteHeader(http.StatusNotFound)

		return
	}

	decodedValue := codec.Decode(metric)
	_, err = io.WriteString(w, decodedValue)
	if err != nil {
		fmt.Printf("Write string, err=%s\n", err)
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
}
