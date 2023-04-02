package api

import (
	"log"
	"memtracker/internal/memtrack/metrics"
	"net/http"
)

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

func (d *DefaultHandler) processRetrieveGauge(metric metrics.Metrics) ([]byte, int) {
	log.Printf("Read to db%v", metric)
	if metric.Delta != nil {
		return []byte{}, http.StatusBadRequest
	} else {
		res, err := d.DB.ReadByParams(metric.MType, metric.ID)
		log.Printf("Retrieve result %v to db %s, %s", metric, res, err)

		if err != nil {
			return []byte{}, http.StatusNotFound
		} else {
			return res, http.StatusOK
		}
	}
}
