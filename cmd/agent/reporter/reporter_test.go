package reporter

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReporter(t *testing.T) {
	gagugeMetrics := map[string]float64{
		"Gauge1": 1.21,
		"Gauge2": 0.0001,
		"Gauge3": 0.000000000001,
	}

	countersMetrics := map[string]int64{
		"Counter1": 20,
	}

	expectedRequests := []string{
		"/update/gauge/Gauge1/1.21",
		"/update/gauge/Gauge2/0.0001",
		"/update/gauge/Gauge3/1E-12",
		"/update/counter/Counter1/20",
	}

	requests := make([]string, 0, len(expectedRequests))

	srvr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			requests = append(requests, r.URL.Path)
		}
	}))

	reporter := New(srvr.URL)

	reporter.Report(gagugeMetrics, countersMetrics)

	for idx := range expectedRequests {
		require.Equal(t, expectedRequests[idx], requests[idx])
	}
}
