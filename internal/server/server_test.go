package server

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kuzhukin/metrics-collector/internal/metric"
	"github.com/kuzhukin/metrics-collector/internal/server/config"
	"github.com/stretchr/testify/require"
)

func TestServerRoutesWithParams(t *testing.T) {
	srvr, err := createServer(&config.Config{EnableLogger: false})
	require.NoError(t, err)

	testFunc := func(t *testing.T, method string, route string) {
		req := httptest.NewRequest(method, route, nil)
		defer req.Body.Close()

		wr := httptest.NewRecorder()
		req.Header.Set("Content-Type", "text/plain")

		srvr.srvr.Handler.ServeHTTP(wr, req)
	}

	t.Run("POST /update", func(t *testing.T) {
		testFunc(t, http.MethodPost, "/update/gauge/metric/100.1")
	})

	t.Run("POST /counter", func(t *testing.T) {
		testFunc(t, http.MethodPost, "/update/counter/metric/1")
	})

	t.Run("GET /gauge", func(t *testing.T) {
		testFunc(t, http.MethodGet, "/value/gauge/metric")
	})

	t.Run("GET /counter", func(t *testing.T) {
		testFunc(t, http.MethodGet, "/value/counter/metric")
	})

	t.Run("GET / - all metrics", func(t *testing.T) {
		testFunc(t, http.MethodGet, "/")
	})
}

func BenchmarkServerRoutesWithParams(b *testing.B) {
	srvr, err := createServer(&config.Config{EnableLogger: false})
	require.NoError(b, err)

	testFunc := func(b *testing.B, method string, route string) {
		b.StopTimer()

		req := httptest.NewRequest(method, route, nil)
		defer req.Body.Close()

		wr := httptest.NewRecorder()
		req.Header.Set("Content-Type", "text/plain")

		b.ReportAllocs()
		b.StartTimer()

		for i := 0; i < b.N; i++ {
			srvr.srvr.Handler.ServeHTTP(wr, req)
		}
	}

	b.ResetTimer()

	b.Run("POST /update", func(b *testing.B) {
		testFunc(b, http.MethodPost, "/update/gauge/metric/100.1")
	})

	b.Run("POST /counter", func(b *testing.B) {
		testFunc(b, http.MethodPost, "/update/counter/metric/1")
	})

	b.Run("GET /gauge", func(b *testing.B) {
		testFunc(b, http.MethodGet, "/value/gauge/metric")
	})

	b.Run("GET /counter", func(b *testing.B) {
		testFunc(b, http.MethodGet, "/value/counter/metric")
	})

	b.Run("GET / - all metrics", func(b *testing.B) {
		testFunc(b, http.MethodGet, "/")
	})
}

func BenchmarkJSONRouter(b *testing.B) {
	srvr, err := createServer(&config.Config{EnableLogger: false})
	require.NoError(b, err)

	testFunc := func(b *testing.B, method string, route string, m *metric.Metric) {
		b.StopTimer()

		data, err := m.Serialize()
		require.NoError(b, err)

		wr := httptest.NewRecorder()

		b.ReportAllocs()
		b.StartTimer()

		for i := 0; i < b.N; i++ {
			req := httptest.NewRequest(method, route, bytes.NewBuffer(data))
			req.Header.Set("Content-Type", "application/json")

			srvr.srvr.Handler.ServeHTTP(wr, req)
		}
	}

	b.ResetTimer()

	gaugeValue := 100.1
	counterValue := int64(1)

	b.Run("POST /update", func(b *testing.B) {
		testFunc(b, http.MethodPost, "/update/", &metric.Metric{
			ID:    "metric",
			Type:  metric.Gauge,
			Value: &gaugeValue,
		})
	})

	b.Run("POST /counter", func(b *testing.B) {
		testFunc(b, http.MethodPost, "/update/", &metric.Metric{
			ID:    "metric",
			Type:  metric.Counter,
			Delta: &counterValue,
		})
	})

	b.Run("GET /gauge", func(b *testing.B) {
		testFunc(b, http.MethodGet, "/value/", &metric.Metric{
			ID:   "metric",
			Type: metric.Gauge,
		})
	})

	b.Run("GET /counter", func(b *testing.B) {
		testFunc(b, http.MethodGet, "/value/", &metric.Metric{
			ID:   "metric",
			Type: metric.Counter,
		})
	})
}

// goos: linux
// goarch: amd64
// pkg: github.com/kuzhukin/metrics-collector/internal/server
// cpu: Intel(R) Core(TM) i7-9700 CPU @ 3.00GHz
// BenchmarkServerRoutesWithParams/POST_/update-8            614010              1843 ns/op             429 B/op          8 allocs/op
// BenchmarkServerRoutesWithParams/POST_/counter-8           694388              1572 ns/op             387 B/op          6 allocs/op
// BenchmarkServerRoutesWithParams/GET_/gauge-8              632457              1766 ns/op             477 B/op          9 allocs/op
// BenchmarkServerRoutesWithParams/GET_/counter-8            694083              1726 ns/op             452 B/op          8 allocs/op
// BenchmarkServerRoutesWithParams/GET_/_-_all_metrics-8     381571              2887 ns/op            2455 B/op         27 allocs/op
// BenchmarkJSONRouter/POST_/update-8                        195804              7950 ns/op            6511 B/op         26 allocs/op
// BenchmarkJSONRouter/POST_/counter-8                       211593              7875 ns/op            6491 B/op         25 allocs/op
// BenchmarkJSONRouter/GET_/gauge-8                          214294              6452 ns/op            6528 B/op         26 allocs/op
// BenchmarkJSONRouter/GET_/counter-8                        221205              6259 ns/op            6523 B/op         26 allocs/op
// PASS
// ok      github.com/kuzhukin/metrics-collector/internal/server   12.964s
