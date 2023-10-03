package handler

import (
	"fmt"
	"io"
	"net/http"

	"github.com/kuzhukin/metrics-collector/cmd/server/codec"
	"github.com/kuzhukin/metrics-collector/internal/shared"
	"github.com/kuzhukin/metrics-collector/internal/storage"
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
	fmt.Println("Get list handler calling...")

	if r.Method != http.MethodGet {
		fmt.Printf("Endpoint %s supports only GET method\n", shared.ValueEndpoint)
		w.WriteHeader(http.StatusMethodNotAllowed)

		return
	}

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)

	metrics := u.storage.List()
	decoded := codec.DecodeMetricsList(metrics)
	_, err := io.WriteString(w, decoded)
	if err != nil {
		fmt.Printf("Write string, err=%s\n", err)
	}
}
