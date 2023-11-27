package controller

import (
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"time"

	"github.com/kuzhukin/metrics-collector/internal/agent/reporter"
	"github.com/kuzhukin/metrics-collector/internal/zlog"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
)

type Controller struct {
	wg             sync.WaitGroup
	polingInterval int
	reportInterval int

	metricsLock    sync.Mutex
	gaugeMetrics   map[string]float64
	counterMetrics map[string]int64

	reporter reporter.Reporter

	done chan struct{}
}

func New(reporter reporter.Reporter, pollingInterval, reportInterval int) *Controller {
	return &Controller{
		polingInterval: pollingInterval,
		reportInterval: reportInterval,
		gaugeMetrics:   make(map[string]float64),
		counterMetrics: make(map[string]int64),
		reporter:       reporter,
		done:           make(chan struct{}),
	}
}

func (c *Controller) Start() {
	zlog.Logger.Infof("Controller started")
	c.start()
	zlog.Logger.Infof("Controller stopped")
}

func (c *Controller) Stop() {
	close(c.done)
}

func (c *Controller) start() {
	c.startCollector()
	c.startReporter()
	c.startGoPsUtilCollector()

	c.wg.Wait()
}

func (c *Controller) startCollector() {
	c.wg.Add(1)
	go func() {
		defer c.wg.Done()

		pollingTicker := time.NewTicker(time.Second * time.Duration(c.polingInterval))
		defer pollingTicker.Stop()

		for {
			select {
			case <-pollingTicker.C:
				c.collectMetrics()
			case <-c.done:
				return
			}
		}
	}()
}

func (c *Controller) startGoPsUtilCollector() {
	c.wg.Add(1)
	go func() {
		defer c.wg.Done()

		pollingTicker := time.NewTicker(time.Second * time.Duration(c.polingInterval))
		defer pollingTicker.Stop()

		for {
			select {
			case <-pollingTicker.C:
				c.collectGoPsUtilMetrics()
			case <-c.done:
				return
			}
		}
	}()
}

func (c *Controller) startReporter() {
	c.wg.Add(1)
	go func() {
		defer c.wg.Done()

		reportTicker := time.NewTicker(time.Second * time.Duration(c.reportInterval))
		defer reportTicker.Stop()

		for {
			select {
			case <-reportTicker.C:
				c.reporter.Report(c.getMetrics())
			case <-c.done:
				return
			}
		}
	}()
}

func (c *Controller) getMetrics() (map[string]float64, map[string]int64) {
	c.metricsLock.Lock()
	defer c.metricsLock.Unlock()

	gauges := make(map[string]float64, len(c.gaugeMetrics))
	counters := make(map[string]int64, len(c.counterMetrics))

	return gauges, counters
}

func (c *Controller) collectMetrics() {

	c.collectGauge()
	c.collectCounter()
}

func (c *Controller) collectGauge() {
	memstats := &runtime.MemStats{}
	runtime.ReadMemStats(memstats)

	c.metricsLock.Lock()
	defer c.metricsLock.Unlock()

	c.addGauge("Alloc", float64(memstats.Alloc))
	c.addGauge("BuckHashSys", float64(memstats.BuckHashSys))
	c.addGauge("Frees", float64(memstats.Frees))
	c.addGauge("GCCPUFraction", memstats.GCCPUFraction)
	c.addGauge("GCSys", float64(memstats.GCSys))
	c.addGauge("HeapAlloc", float64(memstats.HeapAlloc))
	c.addGauge("HeapIdle", float64(memstats.HeapIdle))
	c.addGauge("HeapInuse", float64(memstats.HeapInuse))
	c.addGauge("HeapObjects", float64(memstats.HeapObjects))
	c.addGauge("HeapReleased", float64(memstats.HeapReleased))
	c.addGauge("HeapSys", float64(memstats.HeapSys))
	c.addGauge("LastGC", float64(memstats.LastGC))
	c.addGauge("Lookups", float64(memstats.Lookups))
	c.addGauge("MCacheInuse", float64(memstats.MCacheInuse))
	c.addGauge("MCacheSys", float64(memstats.MCacheSys))
	c.addGauge("MSpanInuse", float64(memstats.MSpanInuse))
	c.addGauge("MSpanSys", float64(memstats.MSpanSys))
	c.addGauge("Mallocs", float64(memstats.Mallocs))
	c.addGauge("NextGC", float64(memstats.NextGC))
	c.addGauge("NumForcedGC", float64(memstats.NumForcedGC))
	c.addGauge("NumGC", float64(memstats.NumGC))
	c.addGauge("OtherSys", float64(memstats.OtherSys))
	c.addGauge("PauseTotalNs", float64(memstats.PauseTotalNs))
	c.addGauge("StackInuse", float64(memstats.StackInuse))
	c.addGauge("StackSys", float64(memstats.StackSys))
	c.addGauge("Sys", float64(memstats.Sys))
	c.addGauge("TotalAlloc", float64(memstats.TotalAlloc))

	// random value
	c.addGauge("RandomValue", genRandom())
}

func (c *Controller) collectGoPsUtilMetrics() {
	v, err := mem.VirtualMemory()
	if err != nil {
		zlog.Logger.Errorf("read gosputil virt mem, err=%s", err)
		return
	}

	utils, err := cpu.Percent(time.Second, true)
	if err != nil {
		zlog.Logger.Errorf("read gosputil virt mem, err=%s", err)
		return
	}

	c.metricsLock.Lock()
	defer c.metricsLock.Unlock()

	c.addGauge("TotalMemory", float64(v.Total))
	c.addGauge("FreeMemory", float64(v.Free))

	for i, percent := range utils {
		c.addGauge(fmt.Sprintf("CPUutilization%d", i), percent)
	}

}

func (c *Controller) addGauge(name string, value float64) {
	c.gaugeMetrics[name] = value
}

func (c *Controller) collectCounter() {
	c.metricsLock.Lock()
	defer c.metricsLock.Unlock()

	c.counterMetrics["PollCount"]++
}

func genRandom() float64 {
	return rand.Float64()
}
