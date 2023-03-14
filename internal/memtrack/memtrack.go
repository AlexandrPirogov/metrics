// Package collects system's metrics
// To see avaible metrics see gauges.go
package memtrack

import (
	"fmt"
	"log"
	"memtracker/internal/memtrack/trackers"
	"net/http"
	"time"
)

// Collects all types of metrics
// Reads and updates metrics
type memtracker struct {
	MetricsContainer trackers.MetricsTracker
}

// Read metrics and send it to given with given http.Client
type httpMemTracker struct {
	Host string
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
func (h httpMemTracker) ReadAndSend(readInterval time.Duration, sendInterval time.Duration) {
	readTicker := time.NewTicker(readInterval)
	sendTicker := time.NewTicker(sendInterval)
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
	go func() {
		for _, metric := range h.MetricsContainer.Metrics {
			metrics := metric.AsMap()
			for k, v := range metrics {
				url := "http://" + h.Host + "/update/" + fmt.Sprintf("%v/%v/%v", metric, k, v)
				log.Printf("Sending metrics to: %s\n", url)
				resp, err := h.client.Post(url, "text/plain", nil)
				if err != nil {
					log.Print(err)
				}
				defer resp.Body.Close()
			}
		}
	}()
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
	return httpMemTracker{
		Host:       host,
		memtracker: memtracker{trackers.New()},
		client:     client}
}
