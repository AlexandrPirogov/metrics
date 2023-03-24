package metrics

import (
	"encoding/json"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMetricSerialize(t *testing.T) {
	deltaValues := []int64{-1, 0, 1}
	valueValues := []float64{-1.1, 0, 1.1}
	data := []struct {
		Name     string
		Expected Metrics
	}{
		//Correct gauge
		{
			Name: "CorrectGauge",
			Expected: Metrics{
				ID:    "1",
				MType: "gauge",
				Delta: &deltaValues[0],
				Value: nil,
			},
		},
		//Correct counters
		{
			Name: "CorrectCounter",
			Expected: Metrics{
				ID:    "1",
				MType: "counter",
				Delta: nil,
				Value: &valueValues[0],
			},
		},
	}

	for _, data := range data {
		t.Run(data.Name, func(t *testing.T) {
			js, err := json.Marshal(data.Expected)
			if err != nil {
				t.Errorf("got error while marshal %v", err)
			}
			log.Println(js)
			unmarshaled := Metrics{}
			err = json.Unmarshal(js, &unmarshaled)
			if err != nil {
				t.Errorf("got error while unmarshal %v", err)
			}
			assert.EqualValues(t, data.Expected, unmarshaled)
		})
	}
}