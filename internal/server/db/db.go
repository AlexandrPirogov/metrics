package db

import (
	"memtracker/internal/server/db/nosql/pidb"
	"strconv"
	"sync"
)

type Storable interface {
	Write(mtype, mname, val string) error
	Read() []byte
	ReadByParams(mtype, mname string) ([]byte, error)
	ReadValueByParams(mtype, mname string) ([]byte, error)
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

// initDB initialize map for MemStorage
func MemStoageDB() pidb.MemStorage {
	var imap = map[string]map[string][]byte{}
	imap["gauge"] = map[string][]byte{}
	imap["counter"] = map[string][]byte{}
	return pidb.MemStorage{Mutex: sync.RWMutex{}, Metrics: imap}
}
