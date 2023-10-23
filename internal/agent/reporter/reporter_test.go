package reporter

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/kuzhukin/metrics-collector/internal/transport"
	"github.com/stretchr/testify/require"
)

func TestReporter(t *testing.T) {
	gagugeMetrics := map[string]float64{
		"Gauge1": 1.21,
		"Gauge2": 0.000000000001,
	}

	countersMetrics := map[string]int64{
		"Counter1": 20,
	}

	expectedMetrics := []transport.Metric{
		{
			ID:    "Gauge1",
			Type:  "gauge",
			Value: 1.21,
		},
		{
			ID:    "Gauge2",
			Type:  "gauge",
			Value: 0.000000000001,
		},
		{
			ID:    "Counter1",
			Type:  "counter",
			Delta: 20,
		},
	}

	requestNumber := 0
	wait := make(chan struct{})

	srvr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestNumber++
		require.Equal(t, http.MethodPost, r.Method)
		require.Equal(t, updateEndpoint, r.URL.Path)

		data, err := io.ReadAll(r.Body)
		require.NoError(t, err)

		m, err := transport.Desirialize(data)
		require.NoError(t, err)

		require.Equal(t, expectedMetrics[requestNumber-1], *m)

		if requestNumber == (len(gagugeMetrics) + len(countersMetrics)) {
			close(wait)
		}
	}))
	defer srvr.Close()

	reporter := New(srvr.URL)

	reporter.Report(gagugeMetrics, countersMetrics)

	select {
	case <-wait:
	case <-time.After(time.Second * 1):
		require.Fail(t, "timeout")
	}
}
