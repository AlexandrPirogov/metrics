package db

import (
	"memtracker/internal/server/db/nosql/pidb"
	"strconv"
)

var DB *pidb.DB = &pidb.DB{
	Documents: make([]pidb.Document, 0),
}

// Saves metric in DB
func Write(mtype, mname, val string) int {
	value, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return -1
	}
	return DB.InsertMetric(mtype, mname, value)
}

// Returns
func Read() string {
	return DB.Select()
}
