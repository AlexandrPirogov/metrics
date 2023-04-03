// Package collects system's metrics
// To see avaible metrics see gauges.go
package memtrack

import (
	"log"
	"memtracker/internal/config/agent"
	"memtracker/internal/memtrack/http/client"
	"memtracker/internal/memtrack/trackers"
	"strconv"
	"time"
)

// Collects all types of metrics
// Reads and updates metrics
type memtracker struct {
	MetricsContainer trackers.MetricsTracker
}

// Read metrics and send it to given with given http.Client
type httpMemTracker struct {
	Host           string
	PollInterval   int
	ReportInterval int
	client         client.Client
	memtracker
}

// ReadAndSend Starts to read metrics
//
// readInterval -- how often read metrics
//
// sendInterval -- how often send metrics to server
//
// WARNING: Race condition appears
func (h httpMemTracker) ReadAndSend() {
	readTicker := time.NewTicker(time.Second * time.Duration(h.PollInterval))
	sendTicker := time.NewTicker(time.Second * time.Duration(h.ReportInterval))
	for {
		//TODO: fix race condition. Read about mutexes in Go
		select {
		case <-readTicker.C:
			h.update()
		case <-sendTicker.C:
			h.send()
		}
	}
}

// Sends metrics to given host
func (h httpMemTracker) send() {
	for _, metric := range h.MetricsContainer.Metrics {
		mapMetrics := metric.AsMap()
		if metric.String() == "gauge" {
			h.client.SendGauges(metric, mapMetrics)
		} else {
			h.client.SendCounter(metric, mapMetrics)
		}
	}
}

// updates values of tracking metrics
func (h httpMemTracker) update() {
	h.MetricsContainer.InvokeTrackers()
}

// NewHttpMemTracker Creates new instance of HttpMemTracker
//
// Pre-cond: Given client instance and host = addr:port
//
// Post-cond: returns new instance of httpMemTracker
func NewHTTPMemTracker(host string) httpMemTracker {
	cfg := agent.ClientCfg
	pollInterval := cfg.PollInterval[:len(cfg.PollInterval)-1]
	poll, err := strconv.Atoi(string(pollInterval))
	if err != nil {
		log.Fatalf("%v", err)
	}
	reportInterval := cfg.ReportInterval[:len(cfg.ReportInterval)-1]
	report, err := strconv.Atoi(string(reportInterval))
	if err != nil {
		log.Fatalf("%v", err)
	}
	return httpMemTracker{
		Host:           cfg.Address,
		PollInterval:   poll,
		ReportInterval: report,
		memtracker:     memtracker{trackers.New()},
		client:         client.NewClient(host, "application/json"),
	}
}
