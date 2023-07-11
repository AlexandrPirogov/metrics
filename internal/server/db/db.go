package db

import (
	"memtracker/internal/server/db/nosql/pidb"
	"strconv"
	"sync"
)

type Storable interface {
	Write(mtype, mname, val string) error
	Read() string
	ReadByParams(mtype, mname string) (string, error)
}

type DB struct {
	Storage pidb.MemStorage
}

// Saves metric in MemStorage
func (d *DB) Write(mtype, mname, val string) error {
	_, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return err
	}
	return d.Storage.InsertMetric(mtype, mname, val)
}

// Returns
func (d *DB) Read() string {
	return d.Storage.Metrics()
}

// ReadByParams reads metrics from storage by given type and name
//
// Pre-cond: given exicting type and string name for metric
//
// Post-cond: returns string of metrics with given name and type
func (d *DB) ReadByParams(mtype, mname string) (string, error) {
	return d.Storage.Select(mtype, mname)
}

// initDB initialize map for MemStorage
func MemStoageDB() pidb.MemStorage {
	var imap = map[string]map[string]pidb.Document{}
	imap["gauge"] = map[string]pidb.Document{}
	imap["counter"] = map[string]pidb.Document{}
	return pidb.MemStorage{Mutex: sync.RWMutex{}, Documents: imap}
}
