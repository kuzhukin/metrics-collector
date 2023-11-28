package controller

import (
	"sync"
	"testing"
	"time"

	"github.com/kuzhukin/metrics-collector/internal/agent/reporter/mockreporter"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

const reportInterval = 10
const pollingInterval = 2

func TestControllerPolling(t *testing.T) {
	mockReporter := mockreporter.NewReporter(t)
	controller := New(mockReporter, pollingInterval, reportInterval)
	require.Len(t, controller.gaugeMetrics, 0)

	go controller.Start()

	// waiting 2 polling intervals
	const pollIntervalsCount = 2
	time.Sleep(time.Second*pollingInterval*pollIntervalsCount + time.Second)

	controller.Stop()

	// waiting for stop
	time.Sleep(time.Second * 1)

	gauges, counters := controller.getMetrics()
	require.Greater(t, len(gauges), 0)
	require.Greater(t, len(counters), 0)

	for _, m := range allGaugeMetrics {
		require.Contains(t, gauges, m)
	}

	for _, m := range allCounterMetrics {
		require.Contains(t, counters, m)
	}

	require.Equal(t, int64(pollIntervalsCount), counters["PollCount"])
}

func TestControllerReporting(t *testing.T) {
	mockReporter := mockreporter.NewReporter(t)
	controller := New(mockReporter, pollingInterval, reportInterval)
	require.Len(t, controller.gaugeMetrics, 0)

	wg := &sync.WaitGroup{}

	wg.Add(1)
	go controller.Start()

	// waiting without reports
	time.Sleep(time.Second * reportInterval / 2)

	mockReporter.On("Report", mock.Anything, mock.Anything)

	// waiting without reports
	time.Sleep(time.Second * reportInterval)

	controller.Stop()
}

var allGaugeMetrics = []string{
	"Alloc",
	"BuckHashSys",
	"Frees",
	"GCCPUFraction",
	"GCSys",
	"HeapAlloc",
	"HeapIdle",
	"HeapInuse",
	"HeapObjects",
	"HeapReleased",
	"HeapSys",
	"LastGC",
	"Lookups",
	"MCacheInuse",
	"MCacheSys",
	"MSpanInuse",
	"MSpanSys",
	"Mallocs",
	"NextGC",
	"NumForcedGC",
	"NumGC",
	"OtherSys",
	"PauseTotalNs",
	"StackInuse",
	"StackSys",
	"Sys",
	"TotalAlloc",
	"RandomValue",
}

var allCounterMetrics = []string{
	"PollCount",
}
