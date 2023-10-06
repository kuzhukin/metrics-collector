package handler

import (
	"fmt"
	"net/http"

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
	fmt.Printf("Get list handler calling, request=%s\n", r.URL.Path)

	if r.Method != http.MethodGet {
		fmt.Printf("Endpoint %s supports only GET method\n", endpoint.ValueEndpoint)
		w.WriteHeader(http.StatusMethodNotAllowed)

		return
	}

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)

	metrics := u.storage.List()
	decoded := codec.DecodeMetricsList(metrics)
	_, err := w.Write([]byte(decoded))
	if err != nil {
		fmt.Printf("Write string, err=%s\n", err)
	}
}
