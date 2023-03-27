package pidb

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"memtracker/internal/memtrack/metrics"
	"strconv"
	"sync"
	"time"
)

// Local db to imitate storage of metrics
type MemStorage struct {
	Mutex   sync.RWMutex
	Metrics map[string]map[string]Metric
}

// Wrapped JSON
type Metric struct {
	Time time.Time
	Name string `json:"id"`
	Type string `json:"type"`
	Val  string
}

func (d Metric) String() string {
	return fmt.Sprintf("Time: %s, Type: %s, Name: %s, Val: %s", d.Time, d.Type, d.Name, d.Val)
}

func (d Metric) Json() []byte {
	bytes := make([]byte, 0)
	var err error
	if d.Type == "counter" {
		val, _ := strconv.ParseInt(d.Val, 10, 64)
		tmp := struct {
			Name  string `json:"id"`    //Metric name
			Type  string `json:"type"`  // Metric type: gauge or counter
			Delta *int64 `json:"delta"` //Metric's val if passing counter
		}{
			Name:  d.Name,
			Type:  d.Type,
			Delta: &val,
		}
		bytes, err = json.Marshal(tmp)
		if err != nil {
			return []byte{}
		}
	} else {
		val, _ := strconv.ParseFloat(d.Val, 64)
		tmp := struct {
			Name string  `json:"id"`    //Metric name
			Type string  `json:"type"`  // Metric type: gauge or counter
			Val  float64 `json:"value"` //Metric's val if passing gauge
		}{
			Name: d.Name,
			Type: d.Type,
			Val:  val,
		}
		bytes, err = json.Marshal(tmp)
		if err != nil {
			return []byte{}
		}
	}
	log.Printf("Marsahled json into %s", bytes)
	return bytes
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
		return elem.Json(), nil
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
			metrics = append(metrics, doc.Json()...)
		}
	}
	return metrics
}

func (p *MemStorage) ReadValueByParams(mtype, mname string) ([]byte, error) {
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
func (p *MemStorage) insertJSON(mtype, name string, val string) {
	p.Mutex.Lock()
	defer p.Mutex.Unlock()
	document := p.newJSON(mtype, name, val)
	p.Metrics[mtype][name] = document
}

// Pre-cond: given mtype, name and val of metric.
// mtype and name should be one of from package metrics.
//
// Post-condition: creats new Metric instance or returns error
func (p *MemStorage) newJSON(mtype, name, val string) Metric {
	if mtype == "counter" {
		if doc, ok := p.Metrics[mtype][name]; ok {
			docVal, _ := strconv.ParseInt(doc.Val, 10, 64)
			valToAdd, _ := strconv.ParseInt(val, 10, 64)
			val = fmt.Sprintf("%d", valToAdd+docVal)
			log.Printf("INcreased val %s", val)
		}
	}
	return Metric{time.Now(), name, mtype, val}
}
