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
	var err error
	js, err := p.Select(mtype, mname)
	if err != nil {
		return []byte{}, err
	}
	var m metrics.Metrics
	err = json.Unmarshal(js, &m)
	if err != nil {
		return []byte{}, err
	}
	if mtype == "counter" {
		return []byte(fmt.Sprintf("%d", *m.Delta)), nil
	}
	return []byte(fmt.Sprintf("%.3f", *m.Value)), nil
}

// InsertMetric creates and insert metrics in MemStorage
//
// Pre-cond: given mtype, name and val of metric.
// mtype and name should be one of from package metrics.
//
// Post-condition: insert opertaion executed.
// Returns 0 if successed. Otherwise means fail
func (p *MemStorage) InsertMetric(mtype, name, val string) ([]byte, error) {
	if metrics.IsMetricCorrect(mtype, name) != nil {
		errMsg := fmt.Sprintf("given not existing metric %s %s\n", mtype, name)
		log.Println(errMsg)
		return []byte{}, errors.New(errMsg)
	}
	bytes, err := p.insertJSON(mtype, name, val)
	if err != nil {
		return []byte{}, err
	}
	return bytes, err
}

func (p *MemStorage) RestoreMetric(mtype, name, val string) ([]byte, error) {
	if metrics.IsMetricCorrect(mtype, name) != nil {
		errMsg := fmt.Sprintf("given not existing metric %s %s\n", mtype, name)
		log.Println(errMsg)
		return []byte{}, errors.New(errMsg)
	}
	bytes, err := p.restoreJSON(mtype, name, val)
	if err != nil {
		return []byte{}, err
	}
	return bytes, err
}

// Creates Metric by given args and insert it to storage
//
// Pre-cond: given correct args for Metrics
//
// Post-condition: insert opertaion executed.
// Returns 0 if successed. Otherwise means fail
func (p *MemStorage) insertJSON(mtype, name string, val string) ([]byte, error) {
	p.Mutex.Lock()
	defer p.Mutex.Unlock()
	document, err := p.newJSON(mtype, name, val, false)
	if err != nil {
		return []byte{}, err
	}
	p.Metrics[mtype][name] = document
	return document, nil
}

func (p *MemStorage) restoreJSON(mtype, name string, val string) ([]byte, error) {
	p.Mutex.Lock()
	defer p.Mutex.Unlock()
	document, err := p.newJSON(mtype, name, val, true)
	if err != nil {
		return []byte{}, err
	}
	p.Metrics[mtype][name] = document
	return document, nil
}

// Pre-cond: given mtype, name and val of metric.
// mtype and name should be one of from package metrics.
//
// Post-condition: creats new Metric instance or returns error
func (p *MemStorage) newJSON(mtype, name, val string, shouldRestore bool) ([]byte, error) {
	if mtype == "counter" {
		return p.counterJSON(mtype, name, val, shouldRestore)
	}
	return p.gagueJSON(mtype, name, val)
}

// counterJSON converts gauge to JSON
func (p *MemStorage) counterJSON(mtype, name, val string, shouldRestore bool) ([]byte, error) {
	var toUpdate metrics.Metrics
	if doc, ok := p.Metrics[mtype][name]; ok {

		err := json.Unmarshal(doc, &toUpdate)
		if err != nil {
			return []byte{}, err
		}
		toAdd, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			return []byte{}, err
		}
		if !shouldRestore {
			delta := *toUpdate.Delta + toAdd
			toUpdate.Delta = &delta
		} else {
			toUpdate.Delta = &toAdd
		}

	} else {
		toUpdate = metrics.Metrics{}
		toAdd, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			return []byte{}, err
		}
		toUpdate.Delta = &toAdd
	}

	toUpdate.ID = name
	toUpdate.MType = mtype

	return json.Marshal(toUpdate)

}

// gagueJSON converts gauge to JSON
func (p *MemStorage) gagueJSON(mtype, name, val string) ([]byte, error) {
	toInsert := metrics.Metrics{
		ID:    name,
		MType: mtype,
		Delta: nil,
	}

	Val, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return []byte{}, err
	}

	toInsert.Value = &Val
	return json.Marshal(toInsert)
}
