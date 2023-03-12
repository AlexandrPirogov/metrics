// Provides wrapped metrics into structs.
//
// There are two types of metrics: gauges and polls
//
// Gauges stands for measuring something.
// In this pakage gauges are used to measure runtime metrics from runtime package from Gp
//
// Counters stands for counting
package metrics

import (
	"reflect"
)

// Metricalbes entities should update own metrics by Read() int method
// Read returns 0 if success otherwise error code
type Metricable interface {
	Read() int
	String() string
}

func IsMetricCorrect(mtype, name string) int {
	var metrics []Metricable = []Metricable{
		&MemStats{},
		&Polls{},
	}
	for _, metric := range metrics {
		if checkFields(metric, mtype, name) == 0 {
			return 0
		}
	}
	return -1
}

func checkFields(metric Metricable, mtype string, name string) int {
	if metric.String() == mtype {
		return 0
	}
	metricType := reflect.TypeOf(metric).Elem()

	for i := 0; i < metricType.NumField(); i++ {
		if name == metricType.Field(i).Name {
			return 0
		}
	}
	return -1
}
