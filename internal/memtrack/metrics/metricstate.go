package metrics

import "memtracker/internal/kernel/tuples"

// Struct to represent Counter metrics as tuple in FP style
type CounterState struct {
	Name  string   `json:"id"`              //Metric name
	Type  string   `json:"type"`            // Metric type: gauge or counter
	Value *counter `json:"delta,omitempty"` //Metric's val if passing counter
}

func (c CounterState) ToTuple() tuples.Tuple {
	tuple := tuples.Tuple{
		Fields: make(map[string]interface{}),
	}
	tuple.SetField("name", c.Name)
	tuple.SetField("type", c.Type)
	tuple.SetField("value", c.Value)
	return tuple
}

// Struct to represent Counter metrics as tuple in FP style
type GaugeState struct {
	Name  string `json:"id"`              //Metric name
	Type  string `json:"type"`            // Metric type: gauge or counter
	Value *gauge `json:"value,omitempty"` //Metric's val if passing counter
}
