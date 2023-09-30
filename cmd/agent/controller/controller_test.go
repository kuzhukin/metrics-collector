package controller

import (
	"sync"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/kuzhukin/metrics-collector/cmd/agent/reporter/mockreporter"
	"github.com/stretchr/testify/require"
)

func TestControllerPolling(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockReporter := mockreporter.NewMockReporter(ctrl)
	controller := New(mockReporter)
	require.Len(t, controller.gaugeMetrics, 0)

	wg := &sync.WaitGroup{}

	wg.Add(1)
	go controller.Start()

	// waiting 2 polling intervals
	const pollIntervalsCount = 2
	time.Sleep(time.Second*pollIntervalSec*pollIntervalsCount + time.Second)

	controller.Stop()

	// waiting for stop
	time.Sleep(time.Second * 1)

	require.Len(t, controller.gaugeMetrics, len(allGaugeMetrics))
	require.Len(t, controller.counterMetrics, len(allCounterMetrics))
	require.Equal(t, int64(pollIntervalsCount), controller.counterMetrics["PollCount"])
}

func TestControllerReporting(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockReporter := mockreporter.NewMockReporter(ctrl)
	controller := New(mockReporter)
	require.Len(t, controller.gaugeMetrics, 0)

	wg := &sync.WaitGroup{}

	wg.Add(1)
	go controller.Start()

	// waiting without reports
	time.Sleep(time.Second * reportInterval / 2)

	mockReporter.EXPECT().Report(gomock.Any(), gomock.Any())

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
