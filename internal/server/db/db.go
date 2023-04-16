package db

import (
	"encoding/json"
	"log"

	"memtracker/internal/config/server"
	"memtracker/internal/kernel"
	"memtracker/internal/kernel/tuples"
	"memtracker/internal/memtrack/metrics"
	"memtracker/internal/server/db/journal"
	"memtracker/internal/server/db/storage/nosql/pidb"
	"sync"
)

type MetricsStorer interface {
	Write(tuple tuples.Tupler) (tuples.Tupler, error)
	Read(condition tuples.Tupler) ([]tuples.Tupler, error)
}

type DB struct {
	Storage   MetricsStorer
	Journaler journal.Journal
}

// initDB initialize map for MemStorage
func MemStorageDB() *pidb.MemStorage {
	var imap = map[string]map[string]tuples.Tupler{}
	imap["gauge"] = map[string]tuples.Tupler{}
	imap["counter"] = map[string]tuples.Tupler{}
	return &pidb.MemStorage{Mutex: sync.RWMutex{}, Metrics: imap}
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
		if err := d.Journaler.Start(); err != nil || server.JournalCfg.StoreFile == "" {
			log.Printf("Can' start journal %v", err)
		}
	}()
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
		tuple := metric.ToTuple()
		d.Journaler.Restored[metric.ID] = tuple
		log.Printf("Restore %v", d.Journaler.Restored[metric.ID])
	}

	for _, tuple := range d.Journaler.Restored {
		//this increased handle time x2
		// TODO refactor
		res, err := kernel.Read(d.Storage, tuple)
		if len(res) == 0 || err != nil {
			kernel.Write(d.Storage, tuple)
		}
	}
}
