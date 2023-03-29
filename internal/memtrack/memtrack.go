// Package collects system's metrics
// To see avaible metrics see gauges.go
package memtrack

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"memtracker/internal/config"
	"memtracker/internal/memtrack/metrics"
	"memtracker/internal/memtrack/trackers"
	"net/http"
	"time"

	"github.com/caarlos0/env/v7"
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
	memtracker
	client http.Client
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
			for k, v := range mapMetrics {
				val := float64(v.(float64))
				toMarsal := metrics.Metrics{
					ID:    k,
					MType: metric.String(),
					Value: &val,
				}
				url := "http://" + h.Host + "/update/"

				js, err := json.Marshal(toMarsal)
				if err != nil {
					log.Printf("%v", err)
					continue
				}
				buffer := bytes.NewBuffer(js)
				resp, err := h.client.Post(url, "application/json", buffer)

				if err != nil {
					log.Print(err)
					continue
				}
				defer resp.Body.Close()

				_, err = io.ReadAll(resp.Body)
				if err != nil {
					log.Printf("error while readall %v", err)
				}
			}
		} else {
			for k, v := range mapMetrics {
				val := float64(v.(float64))
				del := int64(val)
				toMarsal := metrics.Metrics{
					ID:    k,
					MType: metric.String(),
					Delta: &del,
				}
				url := "http://" + h.Host + "/update/"

				js, err := json.Marshal(toMarsal)
				if err != nil {
					log.Printf("%v", err)
					continue
				}
				buffer := bytes.NewBuffer(js)
				resp, err := h.client.Post(url, "application/json", buffer)

				if err != nil {
					log.Print(err)
					continue
				}
				defer resp.Body.Close()

				_, err = io.ReadAll(resp.Body)
				if err != nil {
					log.Printf("error while readall %v", err)
				}
			}
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
func NewHTTPMemTracker(client http.Client, host string) httpMemTracker {
	cfg := config.ClientConfig{}
	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("error while read config %v", err)
	}
	return httpMemTracker{
		Host:           cfg.Address,
		PollInterval:   int(cfg.PollInterval),
		ReportInterval: int(cfg.ReportInterval),
		memtracker:     memtracker{trackers.New()},
		client:         client}
}
