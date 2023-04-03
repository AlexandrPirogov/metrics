package api

import (
	"memtracker/internal/memtrack/metrics"
	"net/http"
)

// processRetrieve retrieve stored metric value depending on the metric's type
//
// Pre-cond: given metric
//
// Post-cond: If success, returns slice of bytes and http status = 200
// otherwise returns empty bite slice and corresponging http status
func (d *DefaultHandler) processRetrieve(metric metrics.Metrics) ([]byte, int) {
	switch {
	case metric.MType == "gauge":
		return d.processRetrieveGauge(metric)
	case metric.MType == "counter":
		return d.processRetrieveCounter(metric)
	default:
		return []byte{}, http.StatusNotImplemented
	}
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
	} else {
		res, err := d.DB.ReadByParams(metric.MType, metric.ID)
		if err != nil {
			return []byte{}, http.StatusNotFound
		} else {
			return res, http.StatusOK
		}
	}
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
	} else {
		res, err := d.DB.ReadByParams(metric.MType, metric.ID)
		if err != nil {
			return []byte{}, http.StatusNotFound
		} else {
			return res, http.StatusOK
		}
	}
}
