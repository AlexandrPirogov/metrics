package db

import (
	"memtracker/internal/server/db/nosql/pidb"
	"net/http"
	"strconv"
)

// MemStorage containes all metrics
var MemStorage *pidb.MemStorage = &pidb.MemStorage{
	Documents: initDB(),
}

// Saves metric in MemStorage
func Write(mtype, mname, val string) error {
	_, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return err
	}
	return MemStorage.InsertMetric(mtype, mname, val)
}

// Returns
func Read() string {
	return MemStorage.Metrics()
}

// ReadByParams reads metrics from storage by given type and name
//
// Pre-cond: given exicting type and string name for metric
//
// Post-cond: returns string of metrics with given name and type
func ReadByParams(mtype, mname string) (string, int) {
	if res, code := MemStorage.Select(mtype, mname); code == 0 {
		return res, http.StatusOK
	}
	return "", http.StatusNotFound
}

// initDB initialize map for MemStorage
func initDB() map[string]map[string]pidb.Document {
	var imap = map[string]map[string]pidb.Document{}
	imap["gauge"] = map[string]pidb.Document{}
	imap["counter"] = map[string]pidb.Document{}
	return imap
}
