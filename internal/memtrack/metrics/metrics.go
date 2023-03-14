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
	"errors"
	"fmt"
	"reflect"
)

// Metricalbes entities should update own metrics by Read() errpr method
// Read returns nil if success otherwise error
type Metricable interface {
	Read() error
	String() string
}

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
	return errors.New("incorrect metric")
}

func checkFields(metric Metricable, mtype string, name string) error {
	if metric.String() == mtype {
		return nil
	}
	metricType := reflect.TypeOf(metric).Elem()

	for i := 0; i < metricType.NumField(); i++ {
		if name == metricType.Field(i).Name {
			return nil
		}
	}
	return errors.New(fmt.Sprintf("Field %s not found in %s", name, mtype))
}
