package metrics

import "memtracker/internal/kernel/tuples"

// Struct to represent Counter metrics as tuple in FP style
type CounterState struct {
	Name  string   `json:"id"`              //Metric name
	Type  string   `json:"type"`            // Metric type: gauge or counter
	Value *counter `json:"delta,omitempty"` //Metric's val if passing counter
}

// ToTuple convers CounterState to tuple.Tuple
//
// Pre-cond:
//
// Post-cond: tuples instance are returned
func (c CounterState) ToTuple() tuples.Tuple {
	tuple := tuples.Tuple{
		Fields: make(map[string]interface{}),
	}
	tuple.SetField("name", c.Name)
	tuple.SetField("type", c.Type)
	tuple.SetField("value", c.Value)
	return tuple
}

// Struct to represent gauge metrics as tuple in FP style
type GaugeState struct {
	Name  string `json:"id"`              //Metric name
	Type  string `json:"type"`            // Metric type: gauge or counter
	Value *gauge `json:"value,omitempty"` //Metric's val if passing counter
}

// ToTuple convers GaugeState to tuple.Tuple
//
// Pre-cond:
//
// Post-cond: tuples instance are returned
func (g GaugeState) ToTuple() tuples.Tuple {
	tuple := tuples.Tuple{
		Fields: make(map[string]interface{}),
	}
	tuple.SetField("name", g.Name)
	tuple.SetField("type", g.Type)
	tuple.SetField("value", g.Value)
	return tuple
}
