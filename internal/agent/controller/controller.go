package controller

import (
	"fmt"
	"math/rand"
	"runtime"
	"time"

	"github.com/kuzhukin/metrics-collector/internal/agent/reporter"
)

const pollIntervalSec = 2
const reportInterval = 10

type Controller struct {
	gaugeMetrics   map[string]float64
	counterMetrics map[string]int64

	reporter reporter.Reporter

	done chan struct{}
}

func New(reporter reporter.Reporter) *Controller {
	return &Controller{
		gaugeMetrics:   make(map[string]float64),
		counterMetrics: make(map[string]int64),
		reporter:       reporter,
		done:           make(chan struct{}),
	}
}

func (c *Controller) Start() {
	c.loop()
}

func (c *Controller) Stop() {
	close(c.done)
}

func (c *Controller) loop() {
	fmt.Printf("Controller started\n")

	pollingTicker := time.NewTicker(time.Second * pollIntervalSec)
	defer pollingTicker.Stop()

	reportTicker := time.NewTicker(time.Second * reportInterval)
	defer reportTicker.Stop()

	for {
		select {
		case <-pollingTicker.C:
			c.collectMetrics()
		case <-reportTicker.C:
			c.reporter.Report(c.gaugeMetrics, c.counterMetrics)
		case <-c.done:
			fmt.Printf("Controller stopped\n")

			return
		}
	}
}

func (c *Controller) collectMetrics() {
	c.collectGauge()
	c.collectCounter()
}

func (c *Controller) collectGauge() {
	memstats := &runtime.MemStats{}
	runtime.ReadMemStats(memstats)

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

func (c *Controller) addGauge(name string, value float64) {
	c.gaugeMetrics[name] = value
}

func (c *Controller) collectCounter() {
	c.counterMetrics["PollCount"]++
}

func genRandom() float64 {
	return rand.Float64()
}
