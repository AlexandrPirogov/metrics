package metrics

type StateTupler interface {
	SetField(key string, value interface{})
	GetField(Key string) interface{}
}

// Struct to represent Counter metrics as tuple in FP style
type CounterState struct {
	Name  string   `json:"id"`              //Metric name
	Type  string   `json:"type"`            // Metric type: gauge or counter
	Value *counter `json:"delta,omitempty"` //Metric's val if passing counter
}

// Struct to represent Counter metrics as tuple in FP style
type GaugeState struct {
	Name  string `json:"id"`              //Metric name
	Type  string `json:"type"`            // Metric type: gauge or counter
	Value *gauge `json:"value,omitempty"` //Metric's val if passing counter
}
