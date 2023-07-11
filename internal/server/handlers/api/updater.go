package api

import (
	"encoding/json"
	"fmt"
	"log"
	"memtracker/internal/config/server"
	"memtracker/internal/crypt"
	"memtracker/internal/kernel"
	"memtracker/internal/kernel/tuples"
	"memtracker/internal/memtrack/metrics"
)

// processUpdate updates metric value depends on metric's type
//
// Pre-cond: given metric
//
// Post-cond: return result or processing metric.
// If success, returns slice of bytes and http status = 200
// otherwise returns empty bite slice and corresponging http status
func (d *DefaultHandler) processUpdate(tupleList tuples.TupleList) (tuples.TupleList, error) {
	res, err := kernel.Write(d.DB.Storage, tupleList)
	if err != nil {
		log.Printf("err while write counter %v", err)
		return tuples.TupleList{}, err
	}

	if server.ServerCfg.Hash != "" {
		res = d.crypt(res)
	}
	return res, nil
}

func (d *DefaultHandler) crypt(tupleList tuples.TupleList) tuples.TupleList {
	res := tuples.TupleList{}
	for tupleList.Next() {
		h := tupleList.Head()
		name := tuples.ExtractString("name", h)
		switch h.(type) {
		case metrics.CounterState:
			val := tuples.ExtractInt64Pointer("value", h)
			counterHash := crypt.Hash(fmt.Sprintf("%s:counter:%d", name, *val), server.ServerCfg.Hash)
			h = h.SetField("hash", counterHash)
		case metrics.GaugeState:
			val := tuples.ExtractFloat64Pointer("value", h)
			gaugeHash := crypt.Hash(fmt.Sprintf("%s:gauge:%f", name, *val), server.ServerCfg.Hash)
			h = h.SetField("hash", gaugeHash)
			log.Printf("gauge hash: %s", tuples.ExtractString("hash", h))
		}
		res = res.Add(h)
		tupleList = tupleList.Tail()
	}
	return res
}

func (d *DefaultHandler) replicate(tupleList tuples.TupleList) {
	for tupleList.Next() {
		h := tupleList.Head()
		record, _ := json.Marshal(h)
		d.DB.Journaler.Write(record)
		tupleList = tupleList.Tail()
	}
}
