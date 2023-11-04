package handler

import (
	"net/http"

	"github.com/kuzhukin/metrics-collector/internal/server/endpoint"
	"github.com/kuzhukin/metrics-collector/internal/server/parser"
	"github.com/kuzhukin/metrics-collector/internal/server/storage"
	"github.com/kuzhukin/metrics-collector/internal/zlog"
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
	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		zlog.Logger.Infof("Endpoint %s supports only GET method", endpoint.ValueEndpoint)
		w.WriteHeader(http.StatusMethodNotAllowed)

		return
	}

	metric, err := u.parser.Parse(r)
	if err != nil {
		zlog.Logger.Warnf("Parse request path=%s, err=%s", r.URL.Path, err)
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	storedMetric, err := u.storage.Get(metric.Kind, metric.Name)
	if err != nil {
		zlog.Logger.Errorf("storage get kind=%s, name=%s err=%s", metric.Kind, metric.Name, err)
		w.WriteHeader(http.StatusNotFound)

		return
	}

	if err := response(w, r, storedMetric); err != nil {
		zlog.Logger.Warnf("response metric=%v, err=%s", *storedMetric, err)
		http.Error(w, "internal error", http.StatusInternalServerError)
	}
}
