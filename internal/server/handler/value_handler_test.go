package handler

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kuzhukin/metrics-collector/internal/metric"
	"github.com/kuzhukin/metrics-collector/internal/server/parser/mockparser"
	"github.com/kuzhukin/metrics-collector/internal/server/storage/mockstorage"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestGetValueHandler(t *testing.T) {
	mockStorage := mockstorage.NewStorage(t)
	mockParser := mockparser.NewRequestParser(t)

	handler := NewValueHandler(mockStorage, mockParser)

	r := httptest.NewRequest(http.MethodGet, "/value/gauge", nil)
	w := httptest.NewRecorder()

	mockStorage.On("Get", mock.Anything, metric.Gauge, fakeMetric.ID).Return(fakeMetric, nil)
	mockParser.On("Parse", r).Return(fakeMetric, nil)

	handler.ServeHTTP(w, r)

	require.Equal(t, http.StatusOK, w.Code)
	require.Equal(t, "text/plain", w.Header().Get("Content-Type"))
}

func TestGetValueJSONHandler(t *testing.T) {
	mockStorage := mockstorage.NewStorage(t)
	mockParser := mockparser.NewRequestParser(t)

	handler := NewValueHandler(mockStorage, mockParser)

	data, err := fakeMetric.Serialize()
	require.NoError(t, err)

	buff := bytes.NewBuffer(data)

	r := httptest.NewRequest(http.MethodGet, "/value", buff)
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	mockStorage.On("Get", mock.Anything, metric.Gauge, fakeMetric.ID).Return(fakeMetric, nil)
	mockParser.On("Parse", r).Return(fakeMetric, nil)

	handler.ServeHTTP(w, r)

	require.Equal(t, http.StatusOK, w.Code)
	require.Equal(t, "application/json", w.Header().Get("Content-Type"))
}
