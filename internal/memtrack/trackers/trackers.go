package trackers

import (
	"memtracker/internal/memtrack/metrics"
)

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

type MetricsTracker struct {
	Metrics []metrics.Metricable
}

func (g *MetricsTracker) InvokeTrackers() int {
	for _, tracker := range g.Metrics {
		code := tracker.Read()
		if code != 0 {
			return code
		}
	}
	return 0
}
