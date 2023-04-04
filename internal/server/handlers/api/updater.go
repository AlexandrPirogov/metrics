package api

import (
	"fmt"
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

	d.DB.Write(metric.MType, metric.ID, fmt.Sprintf("%d", *metric.Delta))
	body, _ := d.DB.ReadByParams(metric.MType, metric.ID)
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

	body, _ := d.DB.Write(metric.MType, metric.ID, fmt.Sprintf("%.11f", *metric.Value))
	return body, http.StatusOK
}
