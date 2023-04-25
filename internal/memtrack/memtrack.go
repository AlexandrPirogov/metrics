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
	// TODO add here job queue
	readTicker := time.NewTicker(time.Second * time.Duration(h.PollInterval))
	sendTicker := time.NewTicker(time.Second * time.Duration(h.ReportInterval))
	for {
		select {
		case <-readTicker.C:
			go h.update()
		case <-sendTicker.C:
			go h.send()
		}
	}
}

// Sends metrics to given host
func (h httpMemTracker) send() {
	// #TODO add here worker pool
	h.client.Send(h.MetricsContainer.Metrics)
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
func NewHTTPMemTracker() httpMemTracker {
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

	client := client.NewClient(cfg.Address, "application/json")
	go client.Listen()
	return httpMemTracker{
		Host:           cfg.Address,
		PollInterval:   poll,
		ReportInterval: report,
		memtracker:     memtracker{trackers.New()},
		client:         client,
	}
}
