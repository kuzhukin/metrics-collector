package handler

import (
	"errors"
	"net/http"

	"github.com/kuzhukin/metrics-collector/internal/server/codec"
	"github.com/kuzhukin/metrics-collector/internal/server/parser"
	"github.com/kuzhukin/metrics-collector/internal/server/storage"
	"github.com/kuzhukin/metrics-collector/internal/zlog"
)

var _ http.Handler = &BatchUpdateHandler{}

// HTTP handler for updating metrics in batch mode
// POST /updates/
// handles metrics batch update
type BatchUpdateHandler struct {
	storage storage.Storage
	parser  parser.BatchRequestParser
}

func NewBatchUpdateHandler(storage storage.Storage, parser parser.BatchRequestParser) *BatchUpdateHandler {
	return &BatchUpdateHandler{
		storage: storage,
		parser:  parser,
	}
}

func (h *BatchUpdateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	metrics, err := h.parser.BatchParse(r)
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

	if err := h.storage.BatchUpdate(r.Context(), metrics); err != nil {
		zlog.Logger.Errorf("batch updater metrics err=%s\n", err)
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusOK)
}
