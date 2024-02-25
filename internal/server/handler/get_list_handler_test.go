package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kuzhukin/metrics-collector/internal/metric"
	"github.com/kuzhukin/metrics-collector/internal/server/storage/mockstorage"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestGetList(t *testing.T) {
	mockStorage := mockstorage.NewStorage(t)
	handler := NewGetListHandler(mockStorage)

	r := httptest.NewRequest(http.MethodGet, fakeURLPath, nil)
	w := httptest.NewRecorder()

	value := 1.1
	mockStorage.On("List", mock.Anything).Return([]*metric.Metric{{ID: "metric", Type: metric.Gauge, Value: &value}}, nil)

	handler.ServeHTTP(w, r)

	require.Equal(t, http.StatusOK, w.Code)
	require.Equal(t, "text/html", w.Header().Get("Content-Type"))
}
