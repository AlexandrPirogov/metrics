package metrics

type StateTupler interface {
	SetField(key string, value interface{})
	GetField(Key string) interface{}
}

type CounterState struct {
	Name  string   `json:"id"`              //Metric name
	Type  string   `json:"type"`            // Metric type: gauge or counter
	Value *counter `json:"delta,omitempty"` //Metric's val if passing counter
}

type GaugeState struct {
	Name  string `json:"id"`            //Metric name
	Type  string `json:"type"`          // Metric type: gauge or counter
	Value *gauge `json:"val,omitempty"` //Metric's val if passing counter
}
