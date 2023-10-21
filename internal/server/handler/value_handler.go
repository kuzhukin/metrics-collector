package handler

import (
	"net/http"

	"github.com/kuzhukin/metrics-collector/internal/server/codec"
	"github.com/kuzhukin/metrics-collector/internal/server/endpoint"
	. "github.com/kuzhukin/metrics-collector/internal/server/log"
	"github.com/kuzhukin/metrics-collector/internal/server/parser"
	"github.com/kuzhukin/metrics-collector/internal/server/storage"
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
		Logger.Warnf("Endpoint %s supports only GET method\n", endpoint.ValueEndpoint)
		w.WriteHeader(http.StatusMethodNotAllowed)

		return
	}

	metric, err := u.parser.Parse(r)
	if err != nil {
		Logger.Warnf("Parse request path=%s, err=%s\n", r.URL.Path, err)
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	metric, err = u.storage.Get(metric.Kind, metric.Name)
	if err != nil {
		Logger.Warnf("Metrics=%v get err=%s\n", metric, err)
		w.WriteHeader(http.StatusNotFound)

		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)

	decodedValue := codec.DecodeValue(metric)
	_, err = w.Write([]byte(decodedValue))
	if err != nil {
		Logger.Errorf("Write string, err=%s\n", err)
	}
}
