package pidb

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"memtracker/internal/memtrack/metrics"
	"strconv"
	"sync"
)

// Local db to imitate storage of metrics
type MemStorage struct {
	Mutex   sync.RWMutex
	Metrics map[string]map[string][]byte
}

// Select returns code and metric in string representation with given name and type
//
// Pre-cond: given mtype, name and val of metric.
// mtype and name should be one of from package metrics.
//
// Post-condition: returns metric in string representation.
// Returns 0 if successed. Otherwise means fail
func (p *MemStorage) Select(mtype, mname string) ([]byte, error) {
	p.Mutex.Lock()
	defer p.Mutex.Unlock()
	res, err := []byte{}, fmt.Errorf("not found")
	if elem, ok := p.Metrics[mtype][mname]; ok {
		return elem, nil
	}
	return res, err
}

// Metrics returns all metrics in string representions
func (p *MemStorage) ReadAllMetrics() []byte {
	p.Mutex.Lock()
	defer p.Mutex.Unlock()
	metrics := make([]byte, 0)
	for _, types := range p.Metrics {
		for _, doc := range types {
			metrics = append(metrics, doc...)
		}
	}
	return metrics
}

func (p *MemStorage) ReadValueByParams(mtype, mname string) ([]byte, error) {
	p.Mutex.Lock()
	defer p.Mutex.Unlock()
	js, err := p.Select(mtype, mname)
	var m metrics.Metrics
	err = json.Unmarshal(js, &m)
	if err != nil {
		return []byte{}, err
	}
	if mtype == "counter" {
		return []byte(fmt.Sprintf("%d", *m.Delta)), nil
	}
	return []byte(fmt.Sprintf("%f", *m.Value)), nil
}

// InsertMetric creates and insert metrics in MemStorage
//
// Pre-cond: given mtype, name and val of metric.
// mtype and name should be one of from package metrics.
//
// Post-condition: insert opertaion executed.
// Returns 0 if successed. Otherwise means fail
func (p *MemStorage) InsertMetric(mtype, name, val string) error {
	if metrics.IsMetricCorrect(mtype, name) != nil {
		errMsg := fmt.Sprintf("given not existing metric %s %s\n", mtype, name)
		log.Println(errMsg)
		return errors.New(errMsg)
	}
	p.insertJSON(mtype, name, val)
	return nil
}

// Creates Metric by given args and insert it to storage
//
// Pre-cond: given correct args for Metrics
//
// Post-condition: insert opertaion executed.
// Returns 0 if successed. Otherwise means fail
func (p *MemStorage) insertJSON(mtype, name string, val string) []byte {
	p.Mutex.Lock()
	defer p.Mutex.Unlock()
	document := p.newJSON(mtype, name, val)
	p.Metrics[mtype][name] = document
	return document
}

// Pre-cond: given mtype, name and val of metric.
// mtype and name should be one of from package metrics.
//
// Post-condition: creats new Metric instance or returns error
func (p *MemStorage) newJSON(mtype, name, val string) []byte {
	if mtype == "counter" {
		if doc, ok := p.Metrics[mtype][name]; ok {
			var toUpdate metrics.Metrics
			err := json.Unmarshal(doc, &toUpdate)
			if err != nil {

			}
			delta := *toUpdate.Delta
			toAdd, _ := strconv.ParseInt(val, 10, 64)
			delta += toAdd
			toUpdate.Delta = &delta
			toUpdate.ID = name
			toUpdate.MType = mtype
			bytes, err := json.Marshal(toUpdate)
			p.Metrics[mtype][name] = bytes
			return bytes
		} else {
			var toUpdate metrics.Metrics
			json.Unmarshal(doc, &toUpdate)
			toAdd, _ := strconv.ParseInt(val, 10, 64)
			toUpdate.Delta = &toAdd
			toUpdate.ID = name
			toUpdate.MType = mtype
			bytes, _ := json.Marshal(toUpdate)
			p.Metrics[mtype][name] = bytes
			return bytes
		}
	}

	toInsert := metrics.Metrics{
		ID:    name,
		MType: mtype,
		Delta: nil,
	}

	Val, _ := strconv.ParseFloat(val, 64)
	toInsert.Value = &Val
	bytes, _ := json.Marshal(toInsert)
	return bytes
}
