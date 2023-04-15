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
	"encoding/json"
	"errors"
	"fmt"
	"memtracker/internal/kernel/tuples"
	"strconv"
	"strings"
)

// Serializable representation of metric
type Metrics struct {
	ID    string   `json:"id"`              //Metric name
	MType string   `json:"type"`            // Metric type: gauge or counter
	Delta *int64   `json:"delta,omitempty"` //Metric's val if passing counter
	Value *float64 `json:"value,omitempty"` //Metric's val if passing gauge
	Hash  string   `json:"hash,omitempty"`
}

func (m Metrics) ToTuple() tuples.Tupler {
	var tuple tuples.Tupler
	switch m.MType {
	case "counter":
		if m.Delta == nil {
			tuple, _ = createCounterState(m.ID, m.MType, "")
			return tuple
		}
		tuple, _ = createCounterState(m.ID, m.MType, fmt.Sprintf("%d", *m.Delta))
	case "gauge":
		if m.Value == nil {
			tuple, _ = createCounterState(m.ID, m.MType, "")
			return tuple
		}
		tuple, _ = createGaugeState(m.ID, m.MType, fmt.Sprintf("%.20f", *m.Value))
	}

	return tuple
}

func (m Metrics) MarshalJSON() ([]byte, error) {
	type MetricStrAlias Metrics
	type MetricAlias Metrics
	del := "none"
	if m.Delta != nil {
		del = fmt.Sprintf("%d", *m.Delta)
	}

	val := "none"
	if m.Value != nil {
		val = fmt.Sprintf("%.20f", *m.Value)
	}

	mAlias := struct {
		MetricStrAlias
		Delta string `json:"sdelta,omitempty"`
		Value string `json:"svalue,omitempty"`
	}{
		MetricStrAlias: (MetricStrAlias)(m),
		Delta:          del,
		Value:          val,
	}

	alias := struct {
		ID    string   `json:"id"`              //Metric name
		MType string   `json:"type"`            // Metric type: gauge or counter
		Delta *int64   `json:"delta,omitempty"` //Metric's val if passing counter
		Value *float64 `json:"value,omitempty"` //Metric's val if passing gauge
		Hash  string   `json:"hash,omitempty"`  //Metric's val if passing gauge
	}{
		ID:    mAlias.ID,
		MType: m.MType,
		Hash:  m.Hash,
	}

	delta, err := strconv.ParseInt(mAlias.Delta, 10, 64)
	if err != nil {
		alias.Delta = nil
	} else {
		alias.Delta = &delta
	}

	Val, err := strconv.ParseFloat(mAlias.Value, 64)
	if err != nil {
		alias.Value = nil
	} else {
		alias.Value = &Val
	}

	return json.Marshal(alias)
}

func (m *Metrics) UnmarshalJSON(data []byte) error {
	type MetricAlias Metrics
	alias := struct {
		*MetricAlias
	}{
		MetricAlias: (*MetricAlias)(m),
	}

	if err := json.Unmarshal(data, &alias); err != nil {
		return err
	}
	return nil
}

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
	if name == "" {
		return fmt.Errorf("name must be not empty")
	}
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

		if strings.EqualFold(name, mname) {
			return nil
		}
	}
	return fmt.Errorf("field %s not found in %s", name, mtype)
}

func FromTuple(t tuples.Tupler) (Metrics, error) {
	m := Metrics{}
	name, ok := t.GetField("name")
	if !ok {
		return m, errors.New("name field required for metrics")
	}
	mtype, ok := t.GetField("type")
	if !ok {
		return m, errors.New("type field required for metrics")
	}

	val, ok := t.GetField("value")
	if !ok {
		return m, errors.New("value field is required for metrics")
	}

	m.ID = name.(string)
	m.MType = mtype.(string)
	switch mtype {
	case "gauge":
		m.Value = val.(*float64)
	case "counter":
		m.Delta = val.(*int64)
	}
	return m, nil
}
