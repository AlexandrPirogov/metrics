package api

import (
	"fmt"
	"log"
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
		body, err := d.DB.ReadByParams(metric.MType, metric.ID)
		if err != nil {
			log.Printf("err while read after write %v metric: %s %s", err, metric.MType, metric.ID)
		}
		return body, http.StatusOK
	}
}

func (d *DefaultHandler) processUpdateGauge(metric metrics.Metrics) ([]byte, int) {
	if metric.Value == nil || metric.Delta != nil {
		return []byte{}, http.StatusBadRequest
	} else {
		d.DB.Write(metric.MType, metric.ID, fmt.Sprintf("%.11f", *metric.Value))
		body, err := d.DB.ReadByParams(metric.MType, metric.ID)
		if err != nil {
			log.Printf("err while read after write %v metric: %s %s", err, metric.MType, metric.ID)
		}
		return body, http.StatusOK
	}
}
