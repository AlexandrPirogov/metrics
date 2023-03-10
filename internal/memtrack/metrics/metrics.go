package metrics

// Metricalbes entities should update own metrics by Read() int method
// Read returns 0 if success otherwise error code
type Metricable interface {
	Read() int
}
