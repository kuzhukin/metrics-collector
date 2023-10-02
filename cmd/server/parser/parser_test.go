package parser

import (
	"context"
	"net/http"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/kuzhukin/metrics-collector/internal/metric"
	"github.com/kuzhukin/metrics-collector/internal/shared"
	"github.com/stretchr/testify/require"
)

func TestParseRequest(t *testing.T) {
	tests := []struct {
		name           string
		metric         map[string]string
		expectedMetric *metric.Metric
		expectedError  error
	}{
		{
			name: "normal gauge",
			metric: map[string]string{
				"kind":  metric.Gauge,
				"name":  "metric",
				"value": "100.1",
			},
			expectedMetric: &metric.Metric{
				Kind:  metric.Gauge,
				Name:  "metric",
				Value: metric.GaugeValue(100.1),
			},
			expectedError: nil,
		},
		{
			name: "normal counter",
			metric: map[string]string{
				"kind":  metric.Counter,
				"name":  "metric",
				"value": "28",
			},
			expectedMetric: &metric.Metric{
				Kind:  metric.Counter,
				Name:  "metric",
				Value: metric.CounterValue(28),
			},
			expectedError: nil,
		},
		{
			name: "without metric name 1",
			metric: map[string]string{
				"kind":  metric.Counter,
				"name":  "",
				"value": "28",
			},
			expectedError: ErrMetricNameIsNotFound,
		},
		{
			name: "without metric name 2",
			metric: map[string]string{
				"kind":  metric.Counter,
				"value": "28",
			}, expectedError: ErrMetricNameIsNotFound,
		},
		{
			name: "bad metric's kind",
			metric: map[string]string{
				"kind":  "bad_kind",
				"name":  "metric",
				"value": "28",
			},
			expectedError: ErrBadMetricKind,
		},
		{
			name: "bad gauge value",
			metric: map[string]string{
				"kind":  metric.Gauge,
				"name":  "metric",
				"value": "aaa",
			},
			expectedError: ErrBadMetricValue,
		},
		{
			name: "bad counter value",
			metric: map[string]string{
				"kind":  metric.Counter,
				"name":  "metric",
				"value": "100.1",
			},
			expectedError: ErrBadMetricValue,
		},
	}

	parser := NewRequestParser()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			metric, err := parser.Parse(prepareRequest(t, test.metric))
			require.ErrorIs(t, err, test.expectedError)
			require.Equal(t, test.expectedMetric, metric)
		})
	}
}

func prepareRequest(t *testing.T, params map[string]string) *http.Request {
	r, err := http.NewRequest(http.MethodPost, shared.UpdateEndpoint, nil)
	require.NoError(t, err)

	routeCtx := chi.NewRouteContext()

	add := func(k, v string, params *chi.RouteParams) *chi.RouteParams {
		params.Keys = append(params.Keys, k)
		params.Values = append(params.Values, v)

		return params
	}

	routeParams := &chi.RouteParams{}
	for k, v := range params {
		routeParams = add(k, v, routeParams)
	}

	routeCtx.URLParams = *routeParams

	ctx := context.WithValue(context.Background(), chi.RouteCtxKey, routeCtx)

	return r.WithContext(ctx)
}
