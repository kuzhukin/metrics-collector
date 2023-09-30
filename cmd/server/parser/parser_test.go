package parser

import (
	"testing"

	"github.com/kuzhukin/metrics-collector/internal/metric"
	"github.com/kuzhukin/metrics-collector/internal/shared"
	"github.com/stretchr/testify/require"
)

func TestParseRequest(t *testing.T) {
	tests := []struct {
		name           string
		request        string
		expectedMetric *metric.Metric
		expectedError  error
	}{
		{
			name:    "normal gauge",
			request: shared.UpdateEndpoint + "gauge/metric/100.1",
			expectedMetric: &metric.Metric{
				Kind:  metric.Gauge,
				Name:  "metric",
				Value: metric.GaugeValue(100.1),
			},
			expectedError: nil,
		},
		{
			name:    "normal counter",
			request: shared.UpdateEndpoint + "counter/metric/28",
			expectedMetric: &metric.Metric{
				Kind:  metric.Counter,
				Name:  "metric",
				Value: metric.CounterValue(28),
			},
			expectedError: nil,
		},
		{
			name:          "without metric name 1",
			request:       shared.UpdateEndpoint + "counter//28",
			expectedError: ErrMetricNameIsNotFound,
		},
		{
			name:          "without metric name 2",
			request:       shared.UpdateEndpoint + "counter/28",
			expectedError: ErrMetricNameIsNotFound,
		},
		{
			name:          "bad metric's kind",
			request:       shared.UpdateEndpoint + "bad_kind/metric/28",
			expectedError: ErrBadMetricKind,
		},
		{
			name:          "bad gauge value",
			request:       shared.UpdateEndpoint + "gauge/metric/aaa",
			expectedError: ErrBadMetricValue,
		},
		{
			name:          "bad counter value",
			request:       shared.UpdateEndpoint + "counter/metric/100.1",
			expectedError: ErrBadMetricValue,
		},
	}

	parser := NewRequestParser()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			metric, err := parser.Parse(test.request)
			require.ErrorIs(t, err, test.expectedError)
			require.Equal(t, test.expectedMetric, metric)
		})
	}
}
