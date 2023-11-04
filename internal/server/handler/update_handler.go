package handler

import (
	"errors"
	"net/http"

	"github.com/kuzhukin/metrics-collector/internal/server/codec"
	"github.com/kuzhukin/metrics-collector/internal/server/endpoint"
	"github.com/kuzhukin/metrics-collector/internal/server/parser"
	"github.com/kuzhukin/metrics-collector/internal/server/storage"
	"github.com/kuzhukin/metrics-collector/internal/zlog"
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
		zlog.Logger.Infof("Endpoint %s supports only POST method", endpoint.UpdateEndpoint)
		w.WriteHeader(http.StatusMethodNotAllowed)

		return
	}

	metric, err := u.parser.Parse(r)
	if err != nil {
		if errors.Is(err, parser.ErrMetricNameIsNotFound) {
			zlog.Logger.Warnf("Metric wasn't found path=%s, err=%s", r.URL.Path, err)
			w.WriteHeader(http.StatusNotFound)
		} else if errors.Is(err, codec.ErrBadMetricValue) {
			zlog.Logger.Warnf("Bad metric value path=%s, err=%s", r.URL.Path, err)
			w.WriteHeader(http.StatusBadRequest)
		} else if errors.Is(err, parser.ErrBadMetricKind) {
			zlog.Logger.Warnf("Bad metric kind path=%s, err=%s", r.URL.Path, err)
			w.WriteHeader(http.StatusBadRequest)
		} else {
			zlog.Logger.Warnf("Parse request path=%s, err=%s", r.URL.Path, err)
			w.WriteHeader(http.StatusBadRequest)
		}

		return
	}

	if err := u.storage.Update(metric); err != nil {
		zlog.Logger.Errorf("Metrics=%v updating err=%s\n", metric, err)
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	if err := response(w, r, metric); err != nil {
		zlog.Logger.Warnf("response metric=%v, err=%s", *metric, err)
		http.Error(w, "internal error", http.StatusInternalServerError)
	}
}
