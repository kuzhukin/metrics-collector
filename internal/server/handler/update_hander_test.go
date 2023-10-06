package handler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/kuzhukin/metrics-collector/internal/server/codec"
	"github.com/kuzhukin/metrics-collector/internal/server/metric"
	"github.com/kuzhukin/metrics-collector/internal/server/parser"
	"github.com/kuzhukin/metrics-collector/internal/server/parser/mockparser"
	"github.com/kuzhukin/metrics-collector/internal/server/storage/mockstorage"
	"github.com/stretchr/testify/require"
)

const fakeURLPath = "/"

var fakeMetric = &metric.Metric{
	Name:  "fake-metric",
	Kind:  metric.Gauge,
	Value: metric.GaugeValue(3.14),
}

func TestBadMethod(t *testing.T) {
	mockStorage := mockstorage.NewStorage(t)
	mockParser := mockparser.NewRequestParser(t)

	handler := NewUpdateHandler(mockStorage, mockParser)

	r := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, r)
	require.Equal(t, http.StatusMethodNotAllowed, w.Code)
}

func TestSuccessUpload(t *testing.T) {
	mockParser := mockparser.NewRequestParser(t)
	mockStorage := mockstorage.NewStorage(t)
	handler := NewUpdateHandler(mockStorage, mockParser)

	r := httptest.NewRequest(http.MethodPost, fakeURLPath, nil)
	w := httptest.NewRecorder()

	mockParser.On("Parse", r).Return(fakeMetric, nil)
	mockStorage.On("Update", fakeMetric).Return(nil)

	handler.ServeHTTP(w, r)

	require.Equal(t, http.StatusOK, w.Code)
	require.Equal(t, "text/plain", w.Header().Get("Content-Type"))
}

func TestParserErrorToStatusCodes(t *testing.T) {
	mockStorage := mockstorage.NewStorage(t)
	mockParser := mockparser.NewRequestParser(t)
	handler := NewUpdateHandler(mockStorage, mockParser)

	tests := []struct {
		name         string
		parserError  error
		expectedCode int
	}{
		{
			name:         "without metric name",
			parserError:  parser.ErrMetricNameIsNotFound,
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "bad metric value",
			parserError:  codec.ErrBadMetricValue,
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "bad metric kind",
			parserError:  parser.ErrBadMetricKind,
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "other error",
			parserError:  errors.New("other error"),
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodPost, fakeURLPath, nil)
			w := httptest.NewRecorder()

			mockParser.On("Parse", r).Return(nil, test.parserError).Once()
			handler.ServeHTTP(w, r)
			require.Equal(t, test.expectedCode, w.Code)
		})
	}
}

func TestStorageErrorToStatusCodes(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mockstorage.NewStorage(t)
	mockParser := mockparser.NewRequestParser(t)

	handler := NewUpdateHandler(mockStorage, mockParser)

	r := httptest.NewRequest(http.MethodPost, fakeURLPath, nil)
	w := httptest.NewRecorder()

	mockParser.On("Parse", r).Return(fakeMetric, nil)
	mockStorage.On("Update", fakeMetric).Return(errors.New("update error"))
	handler.ServeHTTP(w, r)
	require.Equal(t, http.StatusInternalServerError, w.Code)
}
