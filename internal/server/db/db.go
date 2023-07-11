package db

import (
	"encoding/json"
	"fmt"
	"log"
	"memtracker/internal/memtrack/metrics"
	"memtracker/internal/server/db/journal"
	"memtracker/internal/server/db/nosql/pidb"
	"strconv"
	"sync"
)

type Storable interface {
	Write(mtype, mname, val string) ([]byte, error)
	Read() []byte
	ReadByParams(mtype, mname string) ([]byte, error)
	ReadValueByParams(mtype, mname string) ([]byte, error)
}

type DB struct {
	Storage   pidb.MemStorage
	Journaler journal.Journal
}

// Start starts DB with journal
//
// Pre-cond:
//
// Post-cond: db started and ready to work
func (d *DB) Start() {
	d.Journaler = journal.NewJournal()
	bytes, err := d.Journaler.Restore()
	if err == nil {
		log.Printf("restoring db...\n")
		d.restore(bytes)
	} else {
		log.Printf("%v", err)
	}

	go func() {
		if err := d.Journaler.Start(); err != nil {
			log.Printf("Can' start journal %v", err)
		}
	}()
}

// Write writes metric with given type, name and value
//
// Pre-cond: given existing type, non-empty name and correct val
//
// Post-cond: creates metrics from given params and writes to db
// If success returns array of bytes of json'ed metric and nil error
// Otherwise returns empty slice of byte and error
func (d *DB) Write(mtype, mname, val string) ([]byte, error) {
	_, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return []byte{}, err
	}
	bytes, err := d.Storage.InsertMetric(mtype, mname, val)
	go func() { d.Journaler.Write(bytes) }()
	return bytes, err
}

// WriteRestored writes restored from file given type, name and value
// Works like backup
//
// Pre-cond: given existing type, non-empty name and correct val
//
// Post-cond: creates metrics from given params and writes to db
// If success returns array of bytes of json'ed metric and nil error
// Otherwise returns empty slice of byte and error
func (d *DB) WriteRestored(mtype, mname, val string) ([]byte, error) {
	_, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return []byte{}, err
	}
	bytes, err := d.Storage.RestoreMetric(mtype, mname, val)
	// Made it asynch to escape block with unbuffered channel
	go func() { d.Journaler.Write(bytes) }()
	return bytes, err
}

// ReadValueByParams read metric with given params and returns if it exists
//
// Pre-cond: given existing type, non-empty name and correct val
//
// Post-cond: if metric is stored in DB return slice of bytes and nil error
// Otherwise returns empty slice of byte and error
func (d *DB) ReadValueByParams(mtype, mname string) ([]byte, error) {
	return d.Storage.ReadValueByParams(mtype, mname)
}

// Read return all metrics stored in DB
//
// Post-cond: returns all metrics marshaled into slice of byte
func (d *DB) Read() []byte {
	return d.Storage.ReadAllMetrics()
}

// ReadByParams reads metrics from storage by given type and name
//
// Pre-cond: given exicting type and string name for metric
//
// Post-cond: returns string of metrics with given name and type
func (d *DB) ReadByParams(mtype, mname string) ([]byte, error) {
	return d.Storage.Select(mtype, mname)
}

// restore writes restored metric to DB
//
// Pre-cond: given slice of slice of bytes/ marshaled metrics
//
// Post-cond: unmarshal metric and writes it to DB
// If metric has type counter DB will contain last value of counter
func (d *DB) restore(bytes [][]byte) {
	for _, item := range bytes {
		var metric metrics.Metrics
		if err := json.Unmarshal(item, &metric); err != nil {
			continue
		}
		if metric.MType == "counter" {
			d.WriteRestored(metric.MType, metric.ID, fmt.Sprintf("%d", *metric.Delta))
		} else {
			d.WriteRestored(metric.MType, metric.ID, fmt.Sprintf("%.11f", *metric.Value))
		}
	}
}

// initDB initialize map for MemStorage
func MemStoageDB() pidb.MemStorage {
	var imap = map[string]map[string][]byte{}
	imap["gauge"] = map[string][]byte{}
	imap["counter"] = map[string][]byte{}
	return pidb.MemStorage{Mutex: sync.RWMutex{}, Metrics: imap}
}
