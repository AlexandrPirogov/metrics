package api

import (
	"fmt"
	"memtracker/internal/memtrack/metrics"
	"net/http"
)

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

func (d *DefaultHandler) processUpdateCounter(metric metrics.Metrics) ([]byte, int) {
	if metric.Delta == nil || metric.Value != nil {
		return []byte{}, http.StatusBadRequest
	} else {
		d.DB.Write(metric.MType, metric.ID, fmt.Sprintf("%d", *metric.Delta))
		body, _ := d.DB.ReadByParams(metric.MType, metric.ID)
		return body, http.StatusOK
	}
}

func (d *DefaultHandler) processUpdateGauge(metric metrics.Metrics) ([]byte, int) {
	if metric.Value == nil || metric.Delta != nil {
		return []byte{}, http.StatusBadRequest
	} else {
		body, err := d.DB.Write(metric.MType, metric.ID, fmt.Sprintf("%.11f", *metric.Value))
		if err != nil {
		}
		return body, http.StatusOK
	}
}
