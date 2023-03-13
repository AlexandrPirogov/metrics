package pidb

import (
	"fmt"
	"log"
	"memtracker/internal/memtrack/metrics"
	"strconv"
	"time"
)

// Local db to imitate storage of metrics
type MemStorage struct {
	Documents map[string]map[string]Document
}

// Wrapped JSON
type Document struct {
	Time time.Time
	Name string
	Type string
	Val  string
}

func (d Document) String() string {
	return fmt.Sprintf("Time: %s, Type: %s, Name: %s, Val: %s", d.Time, d.Type, d.Name, d.Val)
}

func (p *MemStorage) Select(mtype, mname string) (string, int) {
	res, code := "", -1
	if elem, ok := p.Documents[mtype][mname]; ok {
		return elem.Val, 0
	}
	return res, code
}

func (p *MemStorage) Metrics() string {
	res := ""
	for _, types := range p.Documents {
		for _, doc := range types {
			res += doc.String() + "\n"
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
	p.insertJSON(mtype, name, val)
	return 0
}

// Creates Document by given args and insert it to storage
//
// Pre-cond: given correct args for Metrics
//
// Post-condition: insert opertaion executed.
// Returns 0 if successed. Otherwise means fail
func (p *MemStorage) insertJSON(mtype, name string, val string) {

	document := p.newJSON(mtype, name, val)
	p.Documents[mtype][name] = document

}

// Pre-cond: given mtype, name and val of metric.
// mtype and name should be one of from package metrics.
//
// Post-condition: creats new Document instance or returns error
func (p *MemStorage) newJSON(mtype, name, val string) Document {
	if mtype == "counter" {
		if doc, ok := p.Documents[mtype][name]; ok {
			docVal, _ := strconv.ParseInt(doc.Val, 10, 64)
			valToAdd, _ := strconv.ParseInt(val, 10, 64)
			val = fmt.Sprintf("%d", valToAdd+docVal)
			log.Printf("INcreased val %s", val)
		}
	}
	return Document{time.Now(), name, mtype, val}
}
