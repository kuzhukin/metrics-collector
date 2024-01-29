package handler

import (
	"net/http"

	"github.com/kuzhukin/metrics-collector/internal/server/codec"
	"github.com/kuzhukin/metrics-collector/internal/server/endpoint"
	"github.com/kuzhukin/metrics-collector/internal/server/storage"
	"github.com/kuzhukin/metrics-collector/internal/zlog"
)

var _ http.Handler = &GetListHandler{}

// HTTP handler for getting all metrics in HTML format
// GET /
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
		zlog.Logger.Infof("Endpoint %s supports only GET method", endpoint.ValueEndpoint)
		w.WriteHeader(http.StatusMethodNotAllowed)

		return
	}

	w.Header().Set("Content-Type", "text/html")

	metrics, err := u.storage.List(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

	decoded := codec.DecodeMetricsList(metrics)
	_, err = w.Write([]byte(decoded))
	if err != nil {
		zlog.Logger.Errorf("Write data to response, err=%s", err)
	}
}
