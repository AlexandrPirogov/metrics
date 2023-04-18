package api

import (
	"log"
	"memtracker/internal/config/server"
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
	if server.ServerCfg.Hash != "" {
		tupleList = d.crypt(tupleList)
	}
	body := tuples.MarshalTupleList(tupleList, []byte{})
	return body, http.StatusOK
}
