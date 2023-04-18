package api

import (
	"encoding/json"
	"log"
	"memtracker/internal/kernel"
	"memtracker/internal/kernel/tuples"
	"net/http"
)

// processRetrieve retrieve stored metric value depending on the metric's type
//
// Pre-cond: given metric
//
// Post-cond: If success, returns slice of bytes and http status = 200
// otherwise returns empty bite slice and corresponging http status
func (d *DefaultHandler) processRetrieve(m tuples.Tupler) ([]byte, int) {
	query := m.ToTuple()
	tupleList, err := kernel.Read(d.DB.Storage, query)
	if err != nil {
		return []byte{}, http.StatusBadRequest
	}

	if !tupleList.Next() {
		log.Printf("not found:%v", m)
		return []byte{}, http.StatusNotFound
	}

	body := tuples.MarshalTupleList(tupleList, []byte{})
	return body, http.StatusOK
}

// marshalTuples marshal tuples in slice to slice of bytes
//
// Pre-cond: given slice of tuples
//
// Post-cond: return slice of bytes
func (d *DefaultHandler) marshalTuples(tuples []tuples.Tupler) []byte {
	body := []byte{}
	for _, tuple := range tuples {
		b, err := json.Marshal(tuple)
		if err != nil {
			continue
		}
		body = append(body, b...)
	}
	return body
}
