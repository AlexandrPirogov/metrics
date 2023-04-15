package api

import (
	"encoding/json"
	"fmt"
	"log"
	"memtracker/internal/config/server"
	"memtracker/internal/crypt"
	"memtracker/internal/kernel"
	"memtracker/internal/memtrack/metrics"
	"net/http"
)

// processUpdate updates metric value depends on metric's type
//
// Pre-cond: given metric
//
// Post-cond: return result or processing metric.
// If success, returns slice of bytes and http status = 200
// otherwise returns empty bite slice and corresponging http status
func (d *DefaultHandler) processUpdate(metric metrics.Metrics) ([]byte, int) {
	switch {
	case metric.MType == "gauge":
		return d.processUpdateGauge(metric)
	case metric.MType == "counter":
		return d.processUpdateCounter(metric)
	default:
		return []byte{}, http.StatusNotImplemented
	}
}

// processUpdateCounter updates counter metric
//
// Pre-cond: given counter metric
//
// Post-cond: return result or processing counter  metric.
// If success, returns slice of bytes and http status = 200
// otherwise returns empty bite slice and corresponging http status
func (d *DefaultHandler) processUpdateCounter(metric metrics.Metrics) ([]byte, int) {
	if metric.Delta == nil || metric.Value != nil {
		return []byte{}, http.StatusBadRequest
	}

	key := server.ServerCfg.Hash

	check := crypt.Hash(fmt.Sprintf("%s:counter:%d", metric.ID, *metric.Delta), key)
	if metric.Hash != check {
		log.Printf("Hashes are not equals: \ngot:%s \nhashed:%s", metric.Hash, check)
		return []byte{}, http.StatusBadRequest
	}

	tuple := metric.ToTuple()
	res, err := kernel.Write(d.DB.Storage, tuple)
	if err != nil {
		return []byte{}, http.StatusBadRequest
	}

	body, _ := json.Marshal(res)
	log.Printf("wrote %s", body)
	go func() { d.DB.Journaler.Write(body) }()
	return body, http.StatusOK
}

// processUpdateCounter updates update metric
//
// Pre-cond: given update metric
//
// Post-cond: return result or processing update  metric.
// If success, returns slice of bytes and http status = 200
// otherwise returns empty bite slice and corresponging http status
func (d *DefaultHandler) processUpdateGauge(metric metrics.Metrics) ([]byte, int) {
	if metric.Value == nil || metric.Delta != nil {
		return []byte{}, http.StatusBadRequest
	}

	key := server.ServerCfg.Hash
	check := crypt.Hash(fmt.Sprintf("%s:gauge:%f", metric.ID, *metric.Value), key)
	if metric.Hash != check {
		log.Printf("Hashes are not equals: \ngot:%s \nhashed:%s", metric.Hash, check)
		return []byte{}, http.StatusBadRequest
	}

	tuple := metric.ToTuple()
	res, err := kernel.Write(d.DB.Storage, tuple)
	if err != nil {
		return []byte{}, http.StatusBadGateway
	}
	body, _ := json.Marshal(res)
	go func() { d.DB.Journaler.Write(body) }()
	return body, http.StatusOK
}
