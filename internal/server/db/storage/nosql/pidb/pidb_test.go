package pidb

import (
	"memtracker/internal/kernel"
	"memtracker/internal/kernel/tuples"
	"memtracker/internal/memtrack/metrics"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TODO
// Make Tests better later

// Checks if new created ram db is empty
func TestInit(t *testing.T) {
	db := MemStorage{

		Metrics: initDB(),
	}

	assert.Equal(t, len(db.Metrics), 2, "New created MemStorage must be empty")
}

// Test for saving correct gauges metrics
func TestWriteCorrectGaugesMetrics(t *testing.T) {
	db := MemStorage{
		Metrics: initDB(),
	}

	val := float64(1.1111111)
	var expected = tuples.Tuple{
		Fields: map[string]interface{}{
			"name":  "qwe",
			"type":  "gauge",
			"value": &val,
		},
	}

	actual, err := db.Write(tuples.TupleList{}.Add(expected))

	assert.Nil(t, err)
	assert.EqualValues(t, expected, actual.Head())
}

// Test for saving incorrect gauges metrics
func TestWriteIncorrectGaugesMetrics(t *testing.T) {
	db := MemStorage{

		Metrics: initDB(),
	}
	val := float64(1234567789.123456567)
	var suts = []metrics.GaugeState{
		{
			Name:  "qwe",
			Type:  "",
			Value: &val,
		},
		{
			Name:  "qwe",
			Type:  "counter1",
			Value: &val,
		},
		/*{
			Name:  "qwe",
			Type:  "gauge",
			Value: nil,
		},
		{
			Name:  "qwe",
			Type:  "gauge",
			Value: &val,
		},*/
	}

	for _, sut := range suts {
		t.Run("sut", func(t *testing.T) {

			actual, err := db.Write(tuples.TupleList{}.Add(sut))

			assert.NotNil(t, err)
			assert.NotEqual(t, sut, actual)
		})
	}
}

// Test for saving correct counters metrics

func TestWriteCounterMetrics(t *testing.T) {
	db := MemStorage{

		Metrics: initDB(),
	}

	val := int64(1)
	suts := []metrics.CounterState{
		{
			Name:  "qwe",
			Type:  "counter",
			Value: &val,
		},
		{
			Name:  "qwe",
			Type:  "counter",
			Value: &val,
		},
		{
			Name:  "qwe",
			Type:  "counter",
			Value: &val,
		},
		{
			Name:  "qwe",
			Type:  "counter",
			Value: &val,
		},
	}

	expected := 0
	for _, sut := range suts {
		expected++
		actuals, err := kernel.Write(&db, tuples.TupleList{}.Add(sut))

		for actuals.Next() {
			actualVal, ok := actuals.Head().GetField("value")
			actualV := actualVal.(*int64)
			assert.Nil(t, err)
			assert.True(t, ok)
			assert.EqualValues(t, expected, *actualV)
			actuals = actuals.Tail()
		}
	}

}

// Test for saving incorrect counters metrics

func TestWriteIncorrectCountersMetrics(t *testing.T) {
	db := MemStorage{

		Metrics: initDB(),
	}
	val := int64(1234567789)
	var suts = []metrics.CounterState{
		{
			Name:  "qwe",
			Type:  "",
			Value: &val,
		},
		{
			Name:  "qwe",
			Type:  "counter1",
			Value: &val,
		},
		{
			Name:  "qwe",
			Type:  "gauge",
			Value: nil,
		},
		/*{
			Name:  "qwe",
			Type:  "gauge",
			Value: &val,
		},*/
	}

	for _, sut := range suts {
		t.Run("sut", func(t *testing.T) {

			_, err := db.Write(tuples.TupleList{}.Add(sut))

			assert.NotNil(t, err)
			//	assert.NotEqual(t, sut, actual)
		})
	}
}

func TestReadAll(t *testing.T) {
	db := MemStorage{
		Metrics: filledDB(),
	}

	query := tuples.NewTuple()
	query.SetField("name", "*")
	query.SetField("type", "*")
	actual, err := db.Read(query)

	assert.Nil(t, err)
	assert.True(t, actual.Next())
}

func TestReadAllByParam(t *testing.T) {
	db := MemStorage{Metrics: filledDB()}
	cases := []string{"gauge", "counter"}

	for _, expectedType := range cases {
		t.Run(expectedType, func(t *testing.T) {
			query := tuples.NewTuple()
			query.SetField("name", "*")
			query.SetField("type", expectedType)

			actuals, _ := db.Read(query)

			for actuals.Next() {
				actualType, ok := actuals.Head().GetField("type")
				assert.True(t, ok)
				assert.EqualValues(t, expectedType, actualType.(string))
				actuals = actuals.Tail()
			}
		})
	}
}

func TestReadByParams(t *testing.T) {
	db := MemStorage{Metrics: filledDB()}
	expectedMetrics := db.Metrics

	for _, expectedType := range expectedMetrics {
		for _, expected := range expectedType {
			t.Run("expected", func(t *testing.T) {
				mname, _ := expected.GetField("name")
				mtype, _ := expected.GetField("type")

				query := tuples.NewTuple()
				query.SetField("name", mname.(string))
				query.SetField("type", mtype.(string))

				_, err := db.Read(query)

				assert.Nil(t, err)
				//assert.EqualValues(t, 1, len(actuals))
			})
		}
	}
}

func initDB() map[string]map[string]tuples.Tupler {
	var imap = map[string]map[string]tuples.Tupler{}
	imap["gauge"] = map[string]tuples.Tupler{}
	imap["counter"] = map[string]tuples.Tupler{}
	return imap
}

func filledDB() map[string]map[string]tuples.Tupler {
	var imap = map[string]map[string]tuples.Tupler{}
	imap["gauge"] = map[string]tuples.Tupler{}
	imap["counter"] = map[string]tuples.Tupler{}
	val := float64(1234567789.123456567)
	del := int64(1)
	var suts = []tuples.Tuple{
		{
			Fields: map[string]interface{}{
				"name":  "qwe",
				"type":  "gauge",
				"value": &val,
			},
		},
		{
			Fields: map[string]interface{}{
				"name":  "qwe1",
				"type":  "gauge",
				"value": &val,
			},
		}, {
			Fields: map[string]interface{}{
				"name":  "qwe2",
				"type":  "gauge",
				"value": &val,
			},
		}, {
			Fields: map[string]interface{}{
				"name":  "qwe3",
				"type":  "gauge",
				"value": &val,
			},
		},
		//counters
		{
			Fields: map[string]interface{}{
				"name":  "qwe",
				"type":  "counter",
				"value": &del,
			},
		},
		{
			Fields: map[string]interface{}{
				"name":  "qwe1",
				"type":  "counter",
				"value": &del,
			},
		}, {
			Fields: map[string]interface{}{
				"name":  "qwe2",
				"type":  "counter",
				"value": &del,
			},
		}, {
			Fields: map[string]interface{}{
				"name":  "qwe2",
				"type":  "counter",
				"value": &del,
			},
		},
	}

	for _, sut := range suts {
		name, _ := sut.GetField("name")
		mtype, _ := sut.GetField("type")
		imap[mtype.(string)][name.(string)] = sut
	}

	return imap
}
