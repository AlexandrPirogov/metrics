package pidb

import (
	"memtracker/internal/memtrack/metrics"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Checks if new created ram db is empty
func TestInit(t *testing.T) {
	db := MemStorage{
		make([]Document, 0),
	}

	assert.Equal(t, len(db.Documents), 0, "New created MemStorage must be empty")
}

// Test for saving correct gauges metrics
func TestWriteCorrectGaugesMetrics(t *testing.T) {
	db := MemStorage{
		make([]Document, 0),
	}

	var gauges = metrics.MemStats{}
	gaugesType := reflect.TypeOf(gauges)

	for i := 0; i < gaugesType.NumField(); i++ {
		beforeInsert := len(db.Documents)
		insertStatus := db.InsertMetric(gauges.String(), gaugesType.Field(i).Name, 0)
		afterInsert := len(db.Documents)
		assert.Equal(t, 0, insertStatus, "Can't insert correct gauge metric!\n")
		assert.Greater(t, afterInsert, beforeInsert, "After success insert db size should be increased!")
	}
}

// Test for saving incorrect gauges metrics
func TestWriteIncorrectGaugesMetrics(t *testing.T) {
	db := MemStorage{
		make([]Document, 0),
	}
	var gauges = metrics.MemStats{}
	gaugesType := reflect.TypeOf(gauges)

	for i := 0; i < gaugesType.NumField(); i++ {
		beforeInsert := len(db.Documents)
		modifiedType := " " + gauges.String() + " "
		modifiedName := " " + gaugesType.Field(i).Name + " "
		insertStatus := db.InsertMetric(modifiedType, modifiedName, 0)
		afterInsert := len(db.Documents)
		assert.NotEqual(t, 0, insertStatus, "Can't insert correcet gauge metric!\n")
		assert.Equal(t, afterInsert, beforeInsert, "After success insert db size should be increased!")
	}
}

// Test for saving correct counters metrics
func TestWriteCounterMetrics(t *testing.T) {
	db := MemStorage{
		make([]Document, 0),
	}

	var counters = metrics.Polls{}
	countersType := reflect.TypeOf(counters)

	for i := 0; i < countersType.NumField(); i++ {
		beforeInsert := len(db.Documents)
		insertStatus := db.InsertMetric(counters.String(), countersType.Field(i).Name, 0)
		afterInsert := len(db.Documents)
		assert.Equal(t, 0, insertStatus, "Can't insert correcet gauge metric!\n")
		assert.Greater(t, afterInsert, beforeInsert, "After failed insert db size should be not be modified!")
	}
}

// Test for saving incorrect counters metrics
func TestWriteIncorrectCountersMetrics(t *testing.T) {
	db := MemStorage{
		make([]Document, 0),
	}
	var counters = metrics.Polls{}
	countersType := reflect.TypeOf(counters)

	for i := 0; i < countersType.NumField(); i++ {
		beforeInsert := len(db.Documents)
		modifiedType := " " + counters.String() + " "
		modifiedName := " " + countersType.Field(i).Name + " "
		insertStatus := db.InsertMetric(modifiedType, modifiedName, 0)
		afterInsert := len(db.Documents)
		assert.NotEqual(t, 0, insertStatus, "Can't insert correcet gauge metric!\n")
		assert.Equal(t, afterInsert, beforeInsert, "After success insert db size should be increased!")
	}
}
