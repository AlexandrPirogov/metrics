package pidb

import (
	"encoding/json"
	"fmt"
	"log"
	"memtracker/internal/memtrack/metrics"
	"time"
)

// Local db to imitate storage of metrics
type MemStorage struct {
	Documents map[string]map[string]Document
}

// Wrapped JSON
type Document struct {
	JSON []byte
}

// Metrics to hold
type Metric struct {
	Time time.Time
	Name string
	Type string
	Val  string
}

func (p MemStorage) Select(mtype, mname string) (string, int) {
	res, code := "", -1
	if elem, ok := p.Documents[mtype][mname]; ok {
		var jsonMap map[string]interface{}
		err := json.Unmarshal(elem.JSON, &jsonMap)
		log.Printf("%v\n", jsonMap)
		if err == nil {
			if jsonMap["Name"] == mname && jsonMap["Type"] == mtype {
				res, code = jsonMap["Val"].(string), 0
			}
		}
	}
	return res, code
}

func (p MemStorage) Metrics() string {
	res := ""
	for _, types := range p.Documents {
		for _, doc := range types {
			res += string(doc.JSON) + "\n"
		}
	}
	return res
}

// InsertMetric creates and insert metrics in MemStorage
//
// Pre-cond: given mtype, name and val of metric.
// mtype and name should be one of from package metrics.
//
// Post-condition: insert opertaion executed.
// Returns 0 if successed. Otherwise means fail
func (p *MemStorage) InsertMetric(mtype, name, val string) int {
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
func (p *MemStorage) insertJSON(mtype, name string, val string) int {
	json, err := p.newJSON(mtype, name, val)
	if err != nil {
		log.Print(err)
		return -1
	} else {
		p.Documents[mtype][name] = json
		log.Printf("Successfully saved new Metric %s \n", json)
		return 0
	}
}

// Pre-cond: given mtype, name and val of metric.
// mtype and name should be one of from package metrics.
//
// Post-condition: creats new Document instance or returns error
func (p MemStorage) newJSON(mtype, name, val string) (Document, error) {
	metric := Metric{time.Now(), name, mtype, val}
	js, err := json.Marshal(metric)
	if err != nil {
		log.Print(err)
		return Document{}, fmt.Errorf("%v", err)
	}
	return Document{js}, nil
}
