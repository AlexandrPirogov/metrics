package metrics

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCorrectCreateCounterMetricStateFromCounterMetric(t *testing.T) {
	//Arrange
	delta := int64(1)
	metric := Metrics{
		ID:    "qwe",
		MType: "counter",
		Delta: &delta,
	}

	body, _ := json.Marshal(metric)
	var sut CounterState

	//Act

	err := json.Unmarshal(body, &sut)
	//Assert
	assert.Nil(t, err)
	assert.EqualValues(t, metric.ID, sut.Name)
	assert.EqualValues(t, metric.MType, sut.Type)
	assert.EqualValues(t, metric.Delta, sut.Value)
}

func TestCorrectCreateGaugeMetricStateFromGaugeMetric(t *testing.T) {
	//Arrange
	val := float64(1.11111111)
	metric := Metrics{
		ID:    "qwe",
		MType: "gauge",
		Value: &val,
	}

	body, _ := json.Marshal(metric)
	var sut GaugeState

	//Act

	err := json.Unmarshal(body, &sut)
	//Assert
	assert.Nil(t, err)
	assert.EqualValues(t, metric.ID, sut.Name)
	assert.EqualValues(t, metric.MType, sut.Type)
	assert.EqualValues(t, metric.Value, sut.Value)
}
