package pidb

import (
	"encoding/json"
	"fmt"
	"log"
	"memtracker/internal/memtrack/metrics"
)

// Local db to imitate storage of metrics
type DB struct {
	Documents []Document
}

// Wrapped JSON
type Document struct {
	JSON []byte
}

// Metrics to hold
type Metric struct {
	Name string
	Type string
	Val  float64
}

// InsertMetric creates and insert metrics in DB
//
// Pre-cond: given mtype, name and val of metric.
// mtype and name should be one of from package metrics.
//
// Post-condition: insert opertaion executed.
// Returns 0 if successed. Otherwise means fail
func (p *DB) InsertMetric(mtype, name string, val float64) int {
	if metrics.IsMetricCorrect(mtype, name) != 0 {
		log.Printf("Given not existing metric %s %s\n", mtype, name)
		return -1
	}
	return p.insertJSON(mtype, name, val)
}

// Creates Document by given args and insert it to storage
//
// Pre-cond: given correct args for Metrics
//
// Post-condition: insert opertaion executed.
// Returns 0 if successed. Otherwise means fail
func (p *DB) insertJSON(mtype, name string, val float64) int {
	json, err := p.newJSON(mtype, name, val)
	if err != nil {
		log.Print(err)
		return -1
	} else {
		p.Documents = append(p.Documents, json)
		log.Printf("Successfully saved new Metric. DB holds %d documents \n", len(p.Documents))
		return 0
	}
}

// Pre-cond: given mtype, name and val of metric.
// mtype and name should be one of from package metrics.
//
// Post-condition: creats new Document instance or returns error
func (p DB) newJSON(mtype, name string, val float64) (Document, error) {
	metric := Metric{name, mtype, val}
	js, err := json.Marshal(metric)
	if err != nil {
		log.Print(err)
		return Document{}, fmt.Errorf("%v", err)
	}
	return Document{js}, nil
}
