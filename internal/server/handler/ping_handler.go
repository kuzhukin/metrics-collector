package handler

import (
	"net/http"

	"github.com/kuzhukin/metrics-collector/internal/server/storage/dbstorage"
)

// HTTP handler for checking database connection
// GET /ping
type PingHandler struct {
	db *dbstorage.DBStorage
}

func NewPingHandler(db *dbstorage.DBStorage) *PingHandler {
	return &PingHandler{
		db: db,
	}
}

func (h *PingHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if h.db == nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !h.db.CheckConnection(r.Context()) {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
