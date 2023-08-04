// Package trackers are used to collect all types of measures
//
// Trackers request measures to update own inner metrics
package trackers

import (
	"log"
	"memtracker/internal/metrics"
	"sync"
)

// New Creates new instance of MetricTrackers
//
// Pre-cond:
//
// Post-cond: Creates new instance of MetricsTracker
func New() MetricsTracker {
	polls := metrics.Polls{}
	memStats := metrics.MemStats{}
	opsUtils := metrics.OpsUtil{}
	metricsToCollect := make([]metrics.Metricable, 0)
	metricsToCollect = append(metricsToCollect, &polls)
	metricsToCollect = append(metricsToCollect, &memStats)
	metricsToCollect = append(metricsToCollect, &opsUtils)
	return MetricsTracker{metricsToCollect}
}

type Trackerable interface {
	InvokeTrackers()
}

// MetricsTracker Holds all measures in slice
type MetricsTracker struct {
	Metrics []metrics.Metricable
}

// InvokeTrackers Pre-cond:
//
// Post-cond: requests measures to update own metrics.
// returns 0 if success otherwise we can return error_code
func (g *MetricsTracker) InvokeTrackers() error {
	// #TODO add here grp and errgrp
	var grp sync.WaitGroup
	grp.Add(len(g.Metrics))
	for _, tracker := range g.Metrics {
		go func(tracker metrics.Metricable) {
			defer grp.Done()
			err := tracker.Read()
			if err != nil {
				log.Println(err)
			}
		}(tracker)
	}
	grp.Wait()
	return nil
}
