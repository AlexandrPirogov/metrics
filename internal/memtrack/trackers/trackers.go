// Trackers are used to collect all types of measures
//
// Trackers request measures to update own inner metrics
package trackers

import (
	"memtracker/internal/memtrack/metrics"
)

// Pre-cond:
//
// Post-cond: Creates new instance of MetricsTracker
func New() MetricsTracker {
	polls := metrics.Polls{}
	memStats := metrics.MemStats{}
	metricsToCollect := make([]metrics.Metricable, 0)
	metricsToCollect = append(metricsToCollect, &polls)
	metricsToCollect = append(metricsToCollect, &memStats)
	return MetricsTracker{metricsToCollect}
}

type Trackerable interface {
	InvokeTrackers()
}

// Holds all measures in slice
type MetricsTracker struct {
	Metrics []metrics.Metricable
}

// Pre-cond:
//
// Post-cond: requests measures to update own metrics.
// returns 0 if success otherwise we can return error_code
func (g *MetricsTracker) InvokeTrackers() int {
	for _, tracker := range g.Metrics {
		code := tracker.Read()
		if code != 0 {
			return code
		}
	}
	return 0
}
