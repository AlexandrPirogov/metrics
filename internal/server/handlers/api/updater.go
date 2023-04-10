package api

import (
	"encoding/json"
	"fmt"
	"log"
	"memtracker/internal/config/server"
	"memtracker/internal/crypt"
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
	log.Println("hashinh from server")
	check := crypt.Hash(fmt.Sprintf("%s:counter:%d", metric.ID, *metric.Delta), key)
	if metric.Hash != check {
		log.Printf("Hashes are not equals: \ngot:%s \nhashed:%s", metric.Hash, check)
		return []byte{}, http.StatusBadRequest
	}
	d.DB.Write(metric.MType, metric.ID, fmt.Sprintf("%d", *metric.Delta))
	res, _ := d.DB.ReadByParams(metric.MType, metric.ID)
	var tmp metrics.Metrics
	json.Unmarshal(res, &tmp)
	tmp.Hash = crypt.Hash(fmt.Sprintf("%s:counter:%d", tmp.ID, *tmp.Delta), server.ServerCfg.Hash)
	res, _ = json.Marshal(tmp)

	log.Printf("Body 1after hash: %s", res)
	return res, http.StatusOK
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

	res, _ := d.DB.Write(metric.MType, metric.ID, fmt.Sprintf("%.11f", *metric.Value))

	var tmp metrics.Metrics
	json.Unmarshal(res, &tmp)
	tmp.Hash = crypt.Hash(fmt.Sprintf("%s:gauge:%f", tmp.ID, *tmp.Value), server.ServerCfg.Hash)
	res, _ = json.Marshal(tmp)
	log.Printf("Body 1after hash: %s", res)
	return res, http.StatusOK
}
