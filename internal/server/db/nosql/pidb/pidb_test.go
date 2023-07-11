package pidb

import (
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

	var gauges = metrics.MemStats{}
	metrics := gauges.AsMap()
	for name := range metrics {
		_, err := db.InsertMetric(gauges.String(), name, "0")
		assert.Equal(t, nil, err, "Can't insert correct gauge metric!\n")
	}
}

// Test for saving incorrect gauges metrics
func TestWriteIncorrectGaugesMetrics(t *testing.T) {
	db := MemStorage{

		Metrics: initDB(),
	}
	var gauges = metrics.MemStats{}
	metrics := gauges.AsMap()
	for name := range metrics {
		beforeInsert := len(db.Metrics)
		modifiedType := " " + gauges.String() + " "
		modifiedName := " " + name + " "
		insertStatus, _ := db.InsertMetric(modifiedType, modifiedName, "2")
		afterInsert := len(db.Metrics)
		assert.NotEqual(t, nil, insertStatus, "Can't insert correcet gauge metric!\n")
		assert.Equal(t, afterInsert, beforeInsert, "After success insert db size should be increased!")
	}
}

// Test for saving correct counters metrics
func TestWriteCounterMetrics(t *testing.T) {
	db := MemStorage{

		Metrics: initDB(),
	}

	var counters = metrics.Polls{}
	metrics := counters.AsMap()
	for name := range metrics {
		_, err := db.InsertMetric(counters.String(), name, "0")
		assert.Equal(t, nil, err, "Can't insert correcet gauge metric!\n")
	}
}

// Test for saving incorrect counters metrics
func TestWriteIncorrectCountersMetrics(t *testing.T) {
	db := MemStorage{

		Metrics: initDB(),
	}
	var counters = metrics.Polls{}
	metrics := counters.AsMap()
	for name := range metrics {
		beforeInsert := len(db.Metrics)
		modifiedType := " " + counters.String() + " "
		modifiedName := " " + name + " "
		insertStatus, _ := db.InsertMetric(modifiedType, modifiedName, "2")
		afterInsert := len(db.Metrics)
		assert.NotEqual(t, 0, insertStatus, "Can't insert correcet gauge metric!\n")
		assert.Equal(t, afterInsert, beforeInsert, "After success insert db size should be increased!")
	}
}

func initDB() map[string]map[string][]byte {
	var imap = map[string]map[string][]byte{}
	imap["gauge"] = map[string][]byte{}
	imap["counter"] = map[string][]byte{}
	return imap
}
