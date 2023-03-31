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
		d.Journaler.Start()
	}()
}

// Saves metric in MemStorage
func (d *DB) Write(mtype, mname, val string) ([]byte, error) {
	_, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return []byte{}, err
	}
	bytes, err := d.Storage.InsertMetric(mtype, mname, val)
	go func() { d.Journaler.Write(bytes) }()
	return bytes, err
}

func (d *DB) WriteRestored(mtype, mname, val string) ([]byte, error) {
	_, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return []byte{}, err
	}
	bytes, err := d.Storage.RestoreMetric(mtype, mname, val)
	go func() { d.Journaler.Write(bytes) }()
	return bytes, err
}

func (d *DB) ReadValueByParams(mtype, mname string) ([]byte, error) {
	return d.Storage.ReadValueByParams(mtype, mname)
}

// Returns
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

func (d *DB) StartJournal() error {
	go func() {
		d.Journaler.Start()
	}()
	return nil
}

func (d *DB) restore(bytes [][]byte) {
	for _, item := range bytes {
		var metric metrics.Metrics
		if err := json.Unmarshal(item, &metric); err != nil {
			log.Printf("error while unmasrhal restore %s %v", item, err)
			continue
		}
		if metric.MType == "counter" {
			d.WriteRestored(metric.MType, metric.ID, fmt.Sprintf("%d", *metric.Delta))
		} else {
			d.WriteRestored(metric.MType, metric.ID, fmt.Sprintf("%.11f", *metric.Value))
		}
	}
}

func (d *DB) restoreGauge(mtype string, name string, val string) {
	d.Write(mtype, name, name)
}

// initDB initialize map for MemStorage
func MemStoageDB() pidb.MemStorage {
	var imap = map[string]map[string][]byte{}
	imap["gauge"] = map[string][]byte{}
	imap["counter"] = map[string][]byte{}
	return pidb.MemStorage{Mutex: sync.RWMutex{}, Metrics: imap}
}
