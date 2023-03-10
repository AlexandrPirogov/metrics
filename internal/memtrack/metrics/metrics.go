// Provides wrapped metrics into structs.
//
// There are two types of metrics: gauges and polls
//
// Gauges stands for measuring something.
// In this pakage gauges are used to measure runtime metrics from runtime package from Gp
//
// Counters stands for counting
package metrics

// Metricalbes entities should update own metrics by Read() int method
// Read returns 0 if success otherwise error code
type Metricable interface {
	Read() int
}
