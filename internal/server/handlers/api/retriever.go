package api

import (
	"encoding/json"
	"memtracker/internal/kernel"
	"memtracker/internal/kernel/tuples"
	"memtracker/internal/memtrack/metrics"
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
	res, err := kernel.Read(d.DB.Storage, query)
	if err != nil || len(res) == 0 {
		return []byte{}, http.StatusNotFound
	}

	body := []byte{}
	for _, tuple := range res {
		m, _ := metrics.FromTuple(tuple)
		b, _ := json.Marshal(m)
		body = append(body, b...)
	}
	return body, http.StatusOK
}

// processRetrieve retrieve stored counter metric value
//
// Pre-cond: given counter metric
//
// Post-cond: If success, returns slice of bytes and http status = 200
// otherwise returns empty bite slice and corresponging http status
func (d *DefaultHandler) processRetrieveCounter(metric metrics.Metrics) ([]byte, int) {
	if metric.Delta != nil {
		return []byte{}, http.StatusBadRequest
	}

	query := metric.ToTuple()
	res, err := kernel.Read(d.DB.Storage, query)
	if err != nil {
		return []byte{}, http.StatusNotFound
	}
	body := []byte{}
	for _, tuple := range res {
		m, _ := metrics.FromTuple(tuple)
		b, _ := json.Marshal(m)
		body = append(body, b...)
	}

	return body, http.StatusOK
}

// processRetrieve retrieve stored gauge metric value
//
// Pre-cond: given gauge metric
//
// Post-cond: If success, returns slice of bytes and http status = 200
// otherwise returns empty bite slice and corresponging http status
func (d *DefaultHandler) processRetrieveGauge(metric metrics.Metrics) ([]byte, int) {
	if metric.Delta != nil {
		return []byte{}, http.StatusBadRequest
	}

	query := metric.ToTuple()
	res, err := kernel.Read(d.DB.Storage, query)
	if err != nil {
		return []byte{}, http.StatusNotFound
	}
	body := []byte{}
	for _, tuple := range res {
		m, _ := metrics.FromTuple(tuple)
		b, _ := json.Marshal(m)
		body = append(body, b...)
	}

	return body, http.StatusOK
}
