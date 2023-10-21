package handler

import (
	"net/http"

	"github.com/kuzhukin/metrics-collector/internal/log"
	"github.com/kuzhukin/metrics-collector/internal/server/codec"
	"github.com/kuzhukin/metrics-collector/internal/server/endpoint"
	"github.com/kuzhukin/metrics-collector/internal/server/storage"
)

var _ http.Handler = &GetListHandler{}

type GetListHandler struct {
	storage storage.Storage
}

func NewGetListHandler(storage storage.Storage) *GetListHandler {
	return &GetListHandler{
		storage: storage,
	}
}

func (u *GetListHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		log.Logger.Warnf("Endpoint %s supports only GET method", endpoint.ValueEndpoint)
		w.WriteHeader(http.StatusMethodNotAllowed)

		return
	}

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)

	metrics := u.storage.List()
	decoded := codec.DecodeMetricsList(metrics)
	_, err := w.Write([]byte(decoded))
	if err != nil {
		log.Logger.Errorf("Write data to response, err=%s", err)
	}
}
