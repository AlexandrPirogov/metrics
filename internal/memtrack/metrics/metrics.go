// metrics Provides wrapped metrics into structs.
//
// There are two types of metrics: gauges and polls
//
// Gauges stands for measuring something.
// In this pakage gauges are used to measure runtime metrics from runtime package from Gp
//
// Counters stands for counting
package metrics

import (
	"fmt"
	"strings"
)

// Metricalbes entities should update own metrics by Read() errpr method
// Read returns nil if success otherwise error
type Metricable interface {
	Read() error
	String() string
	AsMap() map[string]interface{}
}

// IsMetricCorrect checks if given type and name for metric is correct
//
// Pre-cond: given correct name and type of metric
//
// Post-cond: return nil if metric is correct, otherwise returns error
func IsMetricCorrect(mtype, name string) error {
	var metrics = []Metricable{
		&MemStats{},
		&Polls{},
	}
	for _, metric := range metrics {
		if checkFields(metric, mtype, name) == nil {
			return nil
		}
	}
	return fmt.Errorf("incorrect metric")
}

// checkFields checks if given type and name exists in given metric
//
// Pre-cond: given correct metric interface name and type of metric
//
// Post-cond: return nil if metric is correct, otherwise returns error
func checkFields(metric Metricable, mtype string, name string) error {
	if metric.String() == mtype {
		return nil
	}

	metrics := metric.AsMap()
	for mname := range metrics {

		if strings.ToLower(name) == strings.ToLower(mname) {
			return nil
		}
	}
	return fmt.Errorf("field %s not found in %s", name, mtype)
}
