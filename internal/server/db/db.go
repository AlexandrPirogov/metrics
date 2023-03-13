package db

import (
	"memtracker/internal/server/db/nosql/pidb"
	"strconv"
)

var MemStorage *pidb.MemStorage = &pidb.MemStorage{
	Documents: make([]pidb.Document, 0),
}

// Saves metric in MemStorage
func Write(mtype, mname, val string) int {
	value, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return -1
	}
	return MemStorage.InsertMetric(mtype, mname, value)
}

// Returns
func Read() string {
	return MemStorage.Select()
}
